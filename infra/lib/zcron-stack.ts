import * as cdk from '@aws-cdk/core'
import * as Lambda from "@aws-cdk/aws-lambda"
import * as Apigateway from "@aws-cdk/aws-apigateway"
import * as Assets from "@aws-cdk/aws-s3-assets"
import * as path from "path"

export class ZcronStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const asset = new Assets.Asset(this, "ZcronLambdaArtifact", { path: path.join(__dirname, "../../bin") })

    const lambdaFunction = new Lambda.Function(this, "ZcronApiHandler", {
      code: Lambda.Code.fromBucket(asset.bucket, asset.s3ObjectKey),
      runtime: Lambda.Runtime.GO_1_X,
      handler: "zcron",
      environment: {
        GIN_MODE: "release"
      }
    })

    new Apigateway.LambdaRestApi(this, "ZcronApiEndpoint", { handler: lambdaFunction })
  }
}
