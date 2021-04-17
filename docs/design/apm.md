# apm

## client
- [demo](https://github.com/meilihao/demo/tree/master/opentelemetry/exporter_otelcol_multi)

## 1. 部署jaeger
1. 根据[Installing Elasticsearch](https://www.elastic.co/guide/en/elasticsearch/reference/7.10/install-elasticsearch.html)安装elasticsearch用于存储jaeger数据
1. 到[jaeger v1.21.0](https://github.com/jaegertracing/jaeger/releases/tag/v1.21.0)下载jaeger-1.21.0-linux-amd64.tar.gz
1. 根据[jaeger v1.21.0中的jaeger-ui git version](https://github.com/jaegertracing/jaeger/tree/v1.21.0)到[其release页面下载对应版本](https://github.com/jaegertracing/jaeger-ui/releases/tag/v1.12.0), 再根据jaeger-ui的README.md `production build`构建出jaeger-query依赖的ui

    jaeger-ui必须使用其release页面下载的版本否则`yarn build`构建可能失败
1. 解压jaeger-1.21.0-linux-amd64.tar.gz到`/opt/jaeger`, 部署jaeger

    ```bash
    # vim run_jaeger.sh
    #!/bin/bash

    export JaegerRoot=${JaegerRoot:-"/opt/jaeger"}
    export ESServerUrls=${ESServerUrls:-"http://localhost:9200"}

    export SPAN_STORAGE_TYPE=elasticsearch # 只有设置SPAN_STORAGE_TYPE=elasticsearch后,collector和query才显示es的配置参数
    nohup ${JaegerRoot}/jaeger-collector 2>&1 > ${JaegerRoot}/collector.log --es.server-urls="${ESServerUrls}" --collector.grpc-server.host-port=":14250" &
    nohup ${JaegerRoot}/jaeger-query 2>&1 > ${JaegerRoot}/query.log --query.static-files=${JaegerRoot}/jaeger-ui/build/ --es.server-urls="${ESServerUrls}" &

    # jaeger-collector --collector.grpc.tls.enabled只能是true. 默认不设置即为grpc insecure=true
    ```

## 2. 部署OpenTelemetry Collector
参考:
- [opentelemetry部署](https://github.com/meilihao/tour_book/blob/master/shell/cmd/suit/opentelemetry.md)

1. 下载[otel-collector_0.18.0_amd64.deb](https://github.com/open-telemetry/opentelemetry-collector/releases).
1. 运行otelcol_linux_amd64

    ```bash
    # cat << EOF > /etc/otel-collector/config.yaml
    # [OpenTelemetry Collector Architecture](https://github.com/open-telemetry/opentelemetry-collector/blob/master/docs/design.md)
    # [Configuring the OpenTelemetry Collector:](https://www.sumologic.com/blog/configure-opentelemetry-collector/)
    extensions:
      health_check:
      pprof:
        endpoint: 0.0.0.0:1777
      zpages:
        endpoint: 0.0.0.0:55679

    receivers:
      otlp:
        protocols:
          grpc:
          # http:

      # # Collect own metrics
      # prometheus:
      #   config:
      #     scrape_configs:
      #       - job_name: 'otel-collector'
      #         scrape_interval: 10s
      #         static_configs:
      #           - targets: ['0.0.0.0:8888']

    processors:
      batch:

    # logging is print in stdout
    exporters:
      logging:
        logLevel: debug
      jaeger:
        endpoint: localhost:14250
        insecure: true # 否则需要设置tls

    service:

      pipelines:

        # traces:
        #   receivers: [otlp]
        #   processors: [batch]
        #   exporters: [logging, jaeger]

        # metrics:
        #   receivers: [otlp]
        #   processors: [batch]
        #   exporters: [logging, prometheus]

        traces:
          receivers: [otlp]
          processors: [batch]
          exporters: [jaeger]
        metrics:
          receivers: [otlp]
          processors: [batch]
          exporters: [logging] # 仅输出到log, 暂时不处理metric
          
      extensions: [health_check, pprof, zpages]
    EOF
    # ./otelcol_linux_amd64 --config otel-config.yaml
    ```