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
      ],
    },
  ] + [
    docker {
      name: 'docker-%s' % name,
      settings+: {
        repo: 'quay.io/signalcd/%s' % name,
        dockerfile: 'cmd/%s/Dockerfile' % name,
      },
    }
    for name in ['api', 'agent']
  ],
}
