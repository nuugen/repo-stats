# repo-stats
A minimal endpoint for retrieving the count of open issues in any publicly listed Github repository.
Included:
- Service written in Golang, exposing the aforementioned endpoint and Dockerized
- CI/CD pipeline implemented with Github Actions.
    * CI is executed automatically every push to `master` branch.
    * CD is only conceptual, deploying to an imaginary K8S cluster, as spec'ed.

## Service
Since the service is not the main focal point of this exercise, Golang has been chosen due to its simplicity for both writers and readers.
Furthermore, [go-github](https://github.com/google/go-github) has been as a Github client, saving some effort in consuming Github's REST API.
For the actual server, the [net/http](https://pkg.go.dev/net/http) and [mux](https://github.com/gorilla/mux) combo is also sufficient.
The example provided on the latter's repository front page was helpful for getting started.

To build and run locally:
```bash
$ docker build -t repo-stats .
$ docker run -p 6699:6699 repo-stats
Listening on http://:6699
```
To test it out:
```bash
$ curl localhost:6699/issue-count/BrandwatchLtd/bcr-api
 {"Status":200,"Message":"","Content":7}
```

## Dockerfile and CI
The a multi-stage Docker build is implemented here, with the primary benefit being a more compact final image containing just the needed final artifacts.
Another benefit is that each intermediate stage can be published as a shareable image.
This could enable a more pleasant experience for onboarding new developers since the long from-scratch compilation can be replaced by a (possibly) quick `docker pull`.

This second point also comes in handy running CI/CD pipelines on a platform with hosted agents, as is in the case of Github Actions. Due to the clean slate nature of every build, caching must be accomplished by pushing intermediate stages.
This precedure is demonstrated in the workflow.

Although not required, it was easy to create a Dockerhub repository for hosting the image. This allows the CI pipeline to be functional on every push `master`.
The registry credentials are stored as secrets in Github and injected into the build environment.

Github Actions was chosen also for its ease of setup in conjunction with the Github repository as well as its YAML syntax.

Reference: https://testdriven.io/blog/faster-ci-builds-with-docker-cache/

## CD to a K8S cluster
This section has not been implemented concretely, but is rather conceptual. Therefore, there several assumptions that we must make with regards to the Kubernetes cluster:
- There should exist a `ServiceAccount` representing the build agent, with the required set of permissions in the form of a `ClusterRole` and of course, a `ClusterRoleBinding` to connect these two objects.
- An `IngressController` should also be present, whether manually installed (Nginx, Traefik, etc..) or can provided by the cloud provider if it's a managed cluster.
- Outbound internet access, so that `repo-stats` can reach Github's API.

The Deployment abstraction is used here as it is the most common way to manage pods.

To expose the application, a Service of type LoadBalancer could have been used, which would map to an external load balancer instance. However, in a realistic setting, it is much more cost-efficient (load balancer instance can be shared) as well as flexible (L7 load balancing) to use an Ingress instead.

The manifest `kubernetes/repo-stats.yml` is definitely not deployable yet, as a few things are missing:
- Implementation-specific Ingress annotations.
- Certificates references for accepting HTTPS requests.
- Functioning healthcheck endpoints in the application.
- Actually defined deployment job in the Github Actions workflow.

Reference: https://kubernetes.io/docs/concepts/services-networking/

## Outlook
Due to the 2-hour constraint, there are things that have been left out that would definitely warrant a follow-up:
- Add the healthcheck endpoints mentioned above.
- Start to write unit test code, the sooner the easier. Include as a stage in Dockerfile.
- Set up TLS. Like unit tests, this tends to also be easier the sooner it is started. To this end [cert-manager](https://cert-manager.io/) is an excellent solution for managing issuance and renewals of TLS certificates in Kubernetes.
