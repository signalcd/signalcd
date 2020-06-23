local pipeline = {
  kind: 'pipeline',
  name: 'build',
  platform: {
    os: 'linux',
    arch: 'amd64',
  },
  steps: [],
};

local golang = {
  name: 'golang',
  image: 'golang:1.14',
  pull: 'always',
  environment: {
    CGO_ENABLED: '0',
    GO111MODULE: 'on',
    GOPROXY: 'https://proxy.golang.org',
  },
};

local node = {
  name: 'node',
  image: 'node:14.4.0',
  pull: 'always',
};

local docker = {
  name: 'docker',
  image: 'plugins/docker',
  settings+: {
    registry: 'quay.io',
    repo: 'quay.io/%s',
    username: {
      from_secret: 'quay_username',
    },
    password: {
      from_secret: 'quay_password',
    },
  },
  when: {
    branch: [
      'master',
    ],
    event: [
      'push',
    ],
  },
};

[
  pipeline {
    steps+: [
      golang {
        name: 'build',
        commands: [
          'make cmd/agent/agent',
          'make cmd/api/api',
        ],
      },
      node {
        name: 'ui',
        commands: [
          'make ui/bundle.js',
          'mkdir -p ./cmd/api/assets',
          'cp ./ui/index.html ./ui/bundle.js ./cmd/api/assets',
        ],
      },
    ] + [
      docker {
        name: 'docker-%s' % name,
        settings+: {
          repo: 'quay.io/signalcd/%s' % name,
          dockerfile: 'cmd/%s/Dockerfile' % name,
          context: './cmd/%s/' % name,
        },
      }
      for name in ['api', 'agent']
    ],
  },

  pipeline {
    name: 'test',

    steps+: [
      golang {
        name: 'test-unit',
        commands: [
          'make test-unit',
        ],
        environment+: {
          CGO_ENABLED: 1,  // for -race
        },
      },
    ],
  },

  pipeline {
    name: 'code-generation',
    steps+: [
      golang {
        name: 'goimports',
        commands: [
          'go get golang.org/x/tools/cmd/goimports',
          'cp $(which goimports) ./goimports',
        ],
      },
    ] + [
      {
        name: '%s' % target,
        image: 'openapitools/openapi-generator-cli:v4.3.1',
        environment: {
          OPENAPI: '/usr/local/bin/docker-entrypoint.sh',
          GOIMPORTS: './goimports',
        },
        commands: [
          'apk add -U git make',
          'make %s --always-make' % target,
          'git diff --exit-code %s' % target,
        ],
      }
      for target in ['api/client/go', 'api/client/javascript', 'api/server/go']
    ],
  },

  pipeline {
    name: 'checks',

    steps+: [
      golang {
        name: 'build-kubernetes-status',
        commands: [
          'make cmd/checks/kubernetes-status/kubernetes-status',
        ],
      },
      docker {
        name: 'docker-kubernetes-status',
        settings+: {
          repo: 'quay.io/signalcd/check-kubernetes-status',
          dockerfile: 'cmd/checks/kubernetes-status/Dockerfile',
          context: 'cmd/checks/kubernetes-status',
        },
      },
      golang {
        name: 'build-http',
        commands: [
          'make cmd/checks/http/http',
        ],
      },
      docker {
        name: 'docker-http',
        settings+: {
          repo: 'quay.io/signalcd/check-http',
          dockerfile: 'cmd/checks/http/Dockerfile',
          context: 'cmd/checks/http',
        },
      },
    ],
  },

  pipeline {
    name: 'plugins',

    steps+: [
      golang {
        name: 'build-drone',
        commands: [
          'make cmd/plugins/drone/drone',
        ],
      },
      docker {
        name: 'docker-drone',
        settings+: {
          repo: 'quay.io/signalcd/drone',
          dockerfile: 'cmd/plugins/drone/Dockerfile',
          context: 'cmd/plugins/drone',
        },
      },
    ],
  },

  pipeline {
    name: 'examples',
    steps+: [
      docker {
        name: 'docker-examples-%s' % name,
        settings+: {
          repo: 'quay.io/signalcd/examples',
          dockerfile: 'examples/%s/Dockerfile' % name,
          context: 'examples/%s' % name,
          tags: [name],
        },
      }
      for name in ['cheese0', 'cheese1', 'cheese2']
    ],
  },
]

