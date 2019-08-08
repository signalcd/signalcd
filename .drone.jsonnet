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
  image: 'golang:1.12',
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

local swagger = {
  name: 'swagger',
  image: 'quay.io/goswagger/swagger:v0.19.0',
};

[
  pipeline {
    steps+: [
      golang {
        name: 'build',
        commands: [
          'make cmd/agent/agent',
          'make cmd/api/api',
          'make cmd/ui/ui',
        ],
      },
      {
        name: 'dart',
        image: 'google/dart:2.3',
        pull: 'always',
        commands: [
          'cd ui',
          'pub get --no-precompile',
          'pub global activate webdev',
          '~/.pub-cache/bin/webdev build',
          'rm -rf build/packages',
          'cp -r build/ ../cmd/ui/assets/',
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
      for name in ['api', 'agent', 'ui']
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
          repo: 'quay.io/signalcd/kubernetes-status',
          dockerfile: 'cmd/checks/kubernetes-status/Dockerfile',
          context: 'cmd/checks/kubernetes-status',
        },
      },
    ],
  },

  pipeline {
    name: 'plugins',

    steps+:[
      golang {
        name: 'build-drone',
        commands: [
          'make cmd/plugins/drone/drone',
        ],
      },
    ],
  },

  pipeline {
    name: 'code-generation',
    steps+: [
      swagger {
        name: 'swagger-apiv1',
        environment: {
          GOSWAGGER: '/usr/bin/swagger',
        },
        commands: [
          'make apiv1',
          'git diff --exit-code',
        ],
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
