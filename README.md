# Kubernetes Admission Webhook

Admission webhooks are HTTP callbacks that receive admission requests and do something with them.
You can define two types of admission webhooks, validating admission Webhook and mutating admission webhook.
With validating admission Webhooks, you may reject requests to enforce custom admission policies.
With mutating admission Webhooks, you may change requests to enforce custom defaults.

This repo currently aims to provide an example for a validating admission Webhook.

## Overview

The Admission webhook requires a `ValidatingWebhookConfiguration` to be created. Once created the Kubernetes API server will send requests to the Webhook based upon the configuration created. 
The configuration specifies a `namespace` & `service` to call back to which will process the webhook, then send a `Allowed` or `Disallowed` to the server. 
In the event the webhook is disallowed, a `Status` response will be added to the request so it is clear why the request was denied.

The example code will reject all deployments that do not have ressource limits set.

## Deploy

Deploy the webhook deployment, service and config:

```bash
$ kubectl apply -f deployment
```