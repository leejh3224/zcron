package main

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"

	"github.com/gorhill/cronexpr"
)

const yyyyMMDDHHmmss = "2006-01-02 15:04:05"

// Returns timezone offset in seconds
func getTimezoneOffsetInSeconds(offset string) int {
	sign := offset[0]
	hh, _ := strconv.Atoi(offset[1:3])
	mm, _ := strconv.Atoi(offset[4:])
	seconds := hh*60*60 + mm*60

	if sign == '-' {
		return -1 * seconds
	}
	return seconds
}

type getNextNSchedulesResponse struct {

	// [supported expressions](https://github.com/gorhill/cronexpr#implementation)
	Expression string `json:"expression"`

	// # of schedules in `nextSchedules`
	Limit int `json:"limit"`

	// schedules adjusted to given `timezoneOffset`
	NextSchedules []string `json:"nextSchedules"`

	// Valid timezone offset of format +-[hh]:[mm]
	TimezoneOffset string `json:"timezoneOffset"`
}

func safeNextN(expr *cronexpr.Expression, n int) (nextSchedules []time.Time) {

	// NextN function might panic in case of accepting invalid expression
	defer func() {
		if err := recover(); err != nil {
			// ignore and just return empty slice
		}
	}()
	nextSchedules = expr.NextN(time.Now().UTC(), uint(n))
	return nextSchedules
}

func getNextNSchedules(c *gin.Context) {
	expression := c.DefaultQuery("expression", "")
	stringLimit := c.DefaultQuery("limit", "5")
	stringTimezoneOffset := c.DefaultQuery("timezoneOffset", "+00:00")

	if len(strings.TrimSpace(expression)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "expression is required"})
		return
	}

	limit, err := strconv.Atoi(stringLimit)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("bad input value `%s` for field `limit`", stringLimit)})
		return
	}

	maxLimit := 50
	if limit > maxLimit {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("`limit` should be smaller than or equal to %d", maxLimit)})
		return
	}

	parsedExpression, err := cronexpr.Parse(expression)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to parse expression, received: `%s`", expression)})
		return
	}

	re := regexp.MustCompile("[+-]\\d{2}:\\d{2}")
	if !re.Match([]byte(stringTimezoneOffset)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to parse timezoneOffset, received: `%s`, wanted: +-[hh]:[mm]", stringTimezoneOffset)})
		return
	}

	nextSchedules := safeNextN(parsedExpression, limit)

	if len(nextSchedules) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to evaluate expression")})
		return
	}

	timezoneOffset := getTimezoneOffsetInSeconds(stringTimezoneOffset)

	schedules := make([]string, 0)

	for _, schedule := range nextSchedules {
		loc := time.FixedZone("customTimezoneOffset", timezoneOffset)
		formatted := schedule.In(loc).Format(yyyyMMDDHHmmss)
		schedules = append(schedules, formatted)
	}

	c.JSON(200, getNextNSchedulesResponse{
		Expression:     expression,
		Limit:          limit,
		NextSchedules:  schedules,
		TimezoneOffset: stringTimezoneOffset,
	})
}

var ginLambda *ginadapter.GinLambda

func init() {
	server := newServer()
	ginLambda = ginadapter.New(server)
}

func newServer() *gin.Engine {
	server := gin.Default()

	server.GET("/next-schedules", getNextNSchedules)

	return server
}

// Handler ...
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
