{
  kind: 'pipeline',
  name: 'build',
  platform: {
    os: 'linux',
    arch: 'amd64',
  },

  local golang = {
    name: 'golang',
    image: 'golang:1.12',
    pull: 'always',
    environment: {
      CGO_ENABLED: '0',
      GO111MODULE: 'on',
      GOPROXY: 'https://proxy.golang.org',
    },
  },

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
        'drone',
      ],
    },
  },

  steps: [
    golang {
      name: 'build',
      commands: [
        'make build',
        // 'make test',
      ],
    },

    // golang {
    //   name: 'generate',
    //   commands: [
    //     'make check-license',
    //     'make generate',
    //     'git diff --exit-code',
    //   ],
    // },

    docker {
      name: 'docker-api',
      settings+: {
        repo: 'quay.io/signalcd/api',
        dockerfile: 'cmd/api/Dockerfile',
      },
    },

    docker {
      name: 'docker-agent',
      settings+: {
        repo: 'quay.io/signalcd/agent',
        dockerfile: 'cmd/agent/Dockerfile',
      },
    },
  ],
}
