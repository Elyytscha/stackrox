# StackRox Kubernetes Security Platform Web Application (UI)

This sub-project contains Web UI (SPA) for StackRox Kubernetes Security Platform.
This project was bootstrapped with [Create React App](https://github.com/facebookincubator/create-react-app).

## Development

If you are developing only StackRox UI, then you don't have to install all the
build tooling described in the parent [README.md](../README.md). Instead, follow
the instructions below.

### Build Tooling

* [Docker](https://www.docker.com/)
* [Node.js](https://nodejs.org/en/) `10.15.3 LTS` or higher (it's highly
recommended to use an LTS version, if you're managing multiple versions of
Node.js on your machine, consider using [nvm](https://github.com/creationix/nvm))
* [Yarn](https://yarnpkg.com/en/)

### Dev Env Setup

_Before starting, make sure you have the above tools installed on your machine
and you've run `yarn install` to download dependencies._

The front end development environment consists of a local static file server
used to serve static UI assets and a remote instance of StackRox for data and
API calls. Set up your environment as follows:

#### Using Local StackRox Deployment and Docker for Mac

_Note: Similar instructions apply when using [Minikube](https://kubernetes.io/docs/setup/minikube/)._

1. **Docker for Mac** - Make sure you have Kubernetes enabled in your Docker for Mac and `kubectl` is
pointing to `docker-desktop` (see [docker docs](https://docs.docker.com/docker-for-mac/#kubernetes)).

1. **Deploy** - Run `yarn deploy-local` (wraps `../deploy/k8s/deploy-local.sh`) to deploy the StackRox software. Make sure that your git working directory is clean and that the branch that you're on has a corresponding tag from CI (see Roxbot comment in a PR). Alternatively, you can check out master before deploying or specify the image tag you want to deploy by setting the `MAIN_IMAGE_TAG` var in your shell.

1. **Start** - Start your local server by running `yarn start`.

_Note: to redeploy a newer version of StackRox, currently the easiest way is by
deleting the whole `stackrox` namespace via `kubectl delete ns stackrox`, and
repeating the steps above._

#### Using a Remote StackRox Deployment

1. **Provision back end cluster** - Navigate to the [Stackrox setup tool](https://setup.rox.systems/). This tool lets you provision a temporary, self destructing cluster in GCloud you will connect to during your development session. Hit the `+` button near the top left. Use the default form settings and provide a "Setup Name" (e.g. `yourname-dev`). Choose the number of hours you would like the cluster to remain active (This should be set to the expected hours of your development session). After you click `run` it may take up to 5 minutes to provision the new cluster. Once the status of your cluster shows as `The cluster is ready`, copy the name of the 'Resource Group' and move on to step 2.

1. **Connect local machine to cluster** - Your local machine needs to be made aware of the cloud cluster you just created. Run `yarn connect [rg-name]` where `[rg-name]` is the name found in **Resource Group** you created in the previous step. This name can be found by going to  https://setup.rox.systems/ and selecting your setup name from the dropdown list.

1. **Deploy StackRox** - Deploy a fresh copy of the StackRox software to your new cluster. During the deployment process, you may be asked for your Dockerhub credentials. In addition to deploying, this command will set up port forwarding from port 8000 to 3000 on your machine.
    * Set up a load balancer by setting the env. variable by running `export LOAD_BALANCER=lb` (optional, if you want to add multiple clusters)
    * Set up persistent storage by setting the env. variable by running `export STORAGE=pvc` (optional, but if you need to bounce central during testing, then your changes will be saved)
    * Run `yarn deploy`

1. **Run local server** - Start your local server.
    * Ensure port forwarding is working by running `yarn forward` (The deploy script tries to do this, but it is flakey.)
    * Start the front-end by running `yarn start`.
    * This will open your web browser to [https://localhost:3000](https://localhost:3000)

_If your machine goes into sleep mode, you may lose the port forwarding set up during the deploy step. If this happens, run `yarn forward` to restart port forwarding._

#### Using an Existing StackRox Deployment with a Local Frontend

If you want to connect your local frontend app to a Stackrox deployment that is already running, you can use one of the following, depending on whether it has a public IP. For both, start by visiting [https://setup.rox.systems/](https://setup.rox.systems/).

* If it has a public IP, find that IP by looking in the **nodes/pods** section in the center-right panel. Copy the **External IP** value. Export that in the `YARN_START_TARGET` env. var and start the front-end by running, `export YARN_START_TARGET=<external_IP>; yarn start`
* If it does not have a public IP, you can steps 2 and 4 from the section above, **Using Remote StackRox Deployment**.
    1. Run `yarn connect [rg-name]` where `[rg-name]` is the name found in the 'Resource Group' found in the existing cluster’s page in Setup.
    1. Run `yarn forward` in one terminal, and `yarn start` in another.

### Testing

#### Unit Tests
Use `yarn test` to run all unit tests and show test coverage.
To run tests and continously watch for changes use `yarn test-watch`.

#### End-to-end Tests (Cypress)

To bring up [Cypress](https://www.cypress.io/) UI use `yarn cypress-open`.
To run all end-to-end tests in a headless mode use `yarn test-e2e-local`.

### Documentation

To start a local server with live-reloading documentation:

```
yarn storybook
```

### IDEs

This project is IDE agnostic. For the best dev experience, it's recommended to
add / configure support for [ESLint](https://eslint.org/) and [Prettier](https://prettier.io/) in the IDE of your choice.

Examples of configuration for some IDEs:

* [Visual Studio Code](https://code.visualstudio.com/): Install plugins [ESLint](https://marketplace.visualstudio.com/items?itemName=dbaeumer.vscode-eslint) and [Prettier](https://marketplace.visualstudio.com/items?itemName=esbenp.prettier-vscode),
then add configuration:

 ```
 "[javascript]": {
    "editor.formatOnSave": true
  },
  "prettier.eslintIntegration": true
```

* [IntelliJ IDEA](https://www.jetbrains.com/idea/) / [WebStorm](https://www.jetbrains.com/webstorm/) / [GoLand](https://www.jetbrains.com/go/): Install and configure [ESLint plugin](https://plugins.jetbrains.com/plugin/7494-eslint). To apply autofixes on file save add [File Watcher](https://www.jetbrains.com/help/idea/using-file-watchers.html) to watch JavaScript files and to run ESLint program `rox/ui/node_modules/.bin/eslint` with arguments `--fix $FilePath$`.

### Browsers

For better development experience it's recommended to use [Google Chrome Browser](https://www.google.com/chrome/) with the following extensions installed:

* [React Developer Tools](https://chrome.google.com/webstore/detail/react-developer-tools/fmkadmapgofadopljbjfkapdkoienihi?hl=en)
* [Redux DevTools](https://chrome.google.com/webstore/detail/redux-devtools/lmhkpmbekcpmknklioeibfkpmmfibljd?hl=en)

### Feature Flags

Just like the backend, we can enable/disable certain features by using environment variables. Create React App will inject environment variables in the app with the prefix `REACT_APP_`. For example, if we want to have a feature flag for Licensing, we could enable the feature by doing a `export REACT_APP_ROX_LICENSE_ENFORCMENT=true` and then starting the app through the usual means. Once a feature is deliverable, we can go ahead and remove the feature flags.
