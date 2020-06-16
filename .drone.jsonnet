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
  image: 'golang:1.13',
  pull: 'always',
  environment: {
    CGO_ENABLED: '0',
    GO111MODULE: 'on',
    GOPROXY: 'https://proxy.golang.org',
  },
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
      {
        name: 'grpc',
        image: 'golang:1.13-alpine',
        environment: {
          GOPROXY: 'https://proxy.golang.org',
        },
        commands: [
          'apk add -U git make protobuf protoc',
          'go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway',
          'go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger',
          'go get -u github.com/golang/protobuf/protoc-gen-go',
          'make signalcd/proto --always-make',
          'git diff --exit-code signalcd/proto',
        ],
      },
      {
        name: 'dart-swagger',
        image: 'openapitools/openapi-generator-cli:v4.3.1',
        environment: {
          SWAGGER: '/usr/local/bin/docker-entrypoint.sh',
        },
        commands: [
          'apk add -U git make',
          'make ui/lib/src/api --always-make',
          'git diff --exit-code ui/lib/src/api',
        ],
      },
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
