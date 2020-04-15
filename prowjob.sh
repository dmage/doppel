#!/bin/sh
curl -sSf 'https://prow.svc.ci.openshift.org/prowjobs.js?omit=annotations,labels,decoration_config,pod_spec'
