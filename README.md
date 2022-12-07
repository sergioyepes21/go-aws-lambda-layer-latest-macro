# go-aws-lambda-layer-latest-macro
A basic Go project that deploys a CloudFormation Macro via AWS SAM. The macro finds the latest Lambda Layer version and returns it on the fragments. 

## Prerequisites
- Golang (go1.19.3+)
- AWS CLI (2.9.0+)
- AWS SAM CLI (1.65.0)

## How to use:
The Macro Function can be invoked by using the `Fn::Transform`:
```yml
MyLambdaFunction:
  Type: AWS::Serverless::Function
  Properties:
    FunctionName: $MyLambdaFunctionName
    ...
    Layers:
      - Fn::Transform:
          Name: $MacroResourceName
          Parameters:
            LayerName: $LambdaLayerArn
```

## Deployment

To deploy the macro into your own AWS Account, you should modify the `samconfig.toml` file with the expected parameters (stack name, S3 deployment bucket, S3 prefix, region and MacroName parameter if needed). 

Later, the following commands should be run:

```shell
$ sam build
$ sam deploy
```

The created resources are:

| LogicalId | Resource Type |
| ------------- | ---- |
| Macro | AWS::CloudFormation::Macro |
| MacroFunction | AWS::Lambda::Function |
| MacroFunctionRole | AWS::IAM::Role |

## Notes

Author: [Sergio Yepes](https://www.linkedin.com/in/sergio-andrés-yepes-joven-41405b174/)

This project was based on Alexis Facques [Python repository](https://github.com/alexisfacques/aws-cloudformation-macro-lambda-latest-layer-version) as an effort to practise my Golang knowledge :).
