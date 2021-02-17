#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from '@aws-cdk/core';
import { ZcronStack } from '../lib/zcron-stack';

const app = new cdk.App();
new ZcronStack(app, 'zcronStack');
