### GRPC

```POST
POST http://127.0.0.1:8080/api/v1/grpc/service    # 获取Service
```

```param
{
    "host":"192.168.1.58:18080",
	"use_tls":false,
	"restart":false
}
```

```Returns
[
    "grpc.reflection.v1alpha.ServerReflection\n",
    "kapi.AnsibleServer\n"
]
```





```POST
POST http://127.0.0.1:8080/api/v1/grpc/function
```

```Param
{
    "host":"192.168.1.58:18080",
    "service": "kapi.AnsibleServer",
	"use_tls":false,
	"restart":false
}
```

```Returns
[
    "kapi.AnsibleServer.CheckConfiguration\n",
    "kapi.AnsibleServer.GenerateYaml\n",
    "kapi.AnsibleServer.RunPlaybook\n",
    "kapi.AnsibleServer.StreamPlaybook\n",
    "kapi.AnsibleServer.StreamRunPlaybook\n"
]
```



```POST
POST http://127.0.0.1:8080/api/v1/grpc/function/param
```

```Param
{
    "host":"192.168.1.58:18080",
    "fun_name": "kapi.AnsibleServer.StreamPlaybook",
	"use_tls":false
}
```

```Returns
{
    "schema": "message PlaybookRequests {\n  string action = 1;\n  repeated .kapi.Node item = 2;\n  .kapi.Config config = 3;\n}",
    "template": "",
    "map_schema": null,
    "map_template": {
        "action": "",
        "config": {
            "clusterName": "",
            "containerNetwork": "",
            "kubeVersion": "",
            "networkMode": "",
            "nfsProvisionerName": "",
            "nfsServer": "",
            "nfsServerPath": ""
        },
        "item": [
            {
                "ip": "",
                "name": "",
                "password": "",
                "port": "",
                "role": ""
            }
        ]
    }
}
```



```POST
POST http://127.0.0.1:8080/api/v1/grpc/function/invoke
```

```Param
{
    "host":"192.168.1.58:18080",
	"fun_name":"kapi.AnsibleServer.GenerateYaml",
	"use_tls":false,
	"body": {
        "item": [
            {
                "ip": "123",
                "name": "124",
                "password": "135",
                "port": "124",
                "role": "153"
            }
        ]
    }
}
```

```Returns
{
    "timer": "2ns",
    "result": "",
    "map_result": {
        "message": "Generate yaml success "
    }
}
```



