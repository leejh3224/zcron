# zcron

`zcron` is an api which offers next schedules for a given cron expression.

It comes handy if you need to adjust schedules to the timezone of your area.

It also supports [AWS cron expression](https://docs.aws.amazon.com/AmazonCloudWatch/latest/events/ScheduledEvents.html) as well.

## Usage

```curl
curl -G --data-urlencode 'expression=0/5 * * * ? *' \
	--data-urlencode 'timezoneOffset=+09:00' \
	--data-urlencode 'limit=50' \
	https://jxx9q1x4ol.execute-api.ap-northeast-2.amazonaws.com/prod/next-schedules | python -m json.tool
```

```bash
# Example Response
{
    "expression": "0/5 * * * ? *",
    "limit": 50,
    "nextSchedules": [
        "2020-11-25 02:40:00",
        "2020-11-25 02:45:00",
        "2020-11-25 02:50:00",
        "2020-11-25 02:55:00",
        "2020-11-25 03:00:00",
        "2020-11-25 03:05:00",
        "2020-11-25 03:10:00",
        "2020-11-25 03:15:00",
        "2020-11-25 03:20:00",
        "2020-11-25 03:25:00"
    ],
    "timezoneOffset": "+09:00"
}
```

## Why I built this

Although AWS has built in cron expression preview, there are some difficulties using it.

- Schedules are based on UTC (you should calculate timzeone offset by yourself)
- Only shows 10 next items

So I tested various cron parser libraries of several languages / online cron generator but almost none of them seemed to work with AWS cron expression well.

It seems [cronexpr](https://github.com/gorhill/cronexpr) seems to be the only one supporting AWS cron expressions, though it's no longer maintained.
