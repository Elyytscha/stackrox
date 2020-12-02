applicationApiVersion: v1beta1

properties:
  ##################################
  #  Required Standard Properties  #
  ##################################
  name:
    type: string
    default: stackrox
    x-google-marketplace:
      type: NAME

  namespace:
    type: string
    default: stackrox
    x-google-marketplace:
      type: NAMESPACE

  #############################
  #  Docker Image Properties  #
  #############################
  main-image:
    type: string
    title: Stackrox image name
    description: Name of Stackrox image to use
    default: gcr.io/stackrox-launcher-project-1/stackrox:$MAIN_IMAGE_TAG
    x-google-marketplace:
      type: IMAGE

  scanner-image:
    type: string
    title: Stackrox Scanner image name
    description: Name of Stackrox scanner image to use
    default: gcr.io/stackrox-launcher-project-1/stackrox/scanner:$SCANNER_IMAGE_TAG
    x-google-marketplace:
      type: IMAGE

  scanner-db-image:
    type: string
    title: Stackrox Scanner DB image name
    description: Name of Stackrox scanner-db image to use
    default: gcr.io/stackrox-launcher-project-1/stackrox/scanner-db:$SCANNER_IMAGE_TAG
    x-google-marketplace:
      type: IMAGE

  ######################
  #  License Property  #
  ######################
  license:
    type: string
    title: (OPTIONAL) Enter your StackRox license key
    description: This is the full license key given to you by StackRox

  ############################
  #  Stackrox.io Properties  #
  ############################
  stackrox-io-username:
    type: string
    title: (OPTIONAL) Enter your stackrox.io username
    description: Used to download artifacts from stackrox.io

  stackrox-io-password:
    type: string
    title: (OPTIONAL) Enter your stackrox.io password
    description: Used to download artifacts from stackrox.io

  ###########################
  #  Networking Properties  #
  ###########################
  network:
    type: string
    title: How do you want to expose StackRox over the network?
    description: This is the method that will be used for exposing StackRox to the network
    default: Load Balancer
    enum:
      - Load Balancer
      - Node Port
      - None

  ########################
  #  Storage Properties  #
  ########################
  pvc-name:
    type: string
    title: What do you want to name the volume?
    description: This is the name that will be given to the persistent volume
    default: stackrox-db

  pvc-storageclass:
    type: string
    title: What storage class do you want to use?
    description: This is the storage class that will be used for the persistent volume
    default: standard

  pvc-size:
    type: integer
    title: How large (in gigabytes) do you want the volume to be?
    description: This is the size in gigabytes that will be allocated for the persistent volume
    default: 100

  ################################
  #  Service Account Properties  #
  ################################
  svcacct:
    type: string
    title: (REQUIRED) Temporary service account to be used when installing StackRox
    x-google-marketplace:
      type: SERVICE_ACCOUNT
      serviceAccount:
        description: This is the temporary service account that will be used to install StackRox
        roles:
        - type: ClusterRole
          rulesType: CUSTOM
          rules:
          - apiGroups: ['*']
            resources: ['*']
            verbs: ['*']

x-google-marketplace:
  clusterConstraints:
    resources:
    - requests:
        memory: 100Mi
        cpu: 100m

required:
- name
- namespace
- main-image
- scanner-image
- scanner-db-image
- network
- pvc-name
- pvc-storageclass
- pvc-size
- svcacct
