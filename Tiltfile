k8s_resource('caddy', port_forwards=6060)
k8s_yaml('development/caddy.yaml')

k8s_resource('signalcd-api', port_forwards=6660)
k8s_yaml('deployment/kubernetes/signalcd-api.yaml')
docker_build('quay.io/signalcd/api', '.', dockerfile='cmd/api/Dockerfile.dev',
  live_update=[
    sync('.', '/go/src/github.com/signalcd/signalcd'),
    run('go install -v /go/src/github.com/signalcd/signalcd/cmd/api'),
    restart_container(),
  ]
)

k8s_resource('signalcd-ui', port_forwards=6662)
k8s_yaml('deployment/kubernetes/signalcd-ui.yaml')
docker_build('quay.io/signalcd/ui', '.', dockerfile='cmd/ui/Dockerfile.dev',
  live_update=[
    fall_back_on(['ui/pubspec.lock', 'ui/pubspec.yaml']),
    sync('ui/', '/app'),
  ]
)
