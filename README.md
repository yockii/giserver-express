# Giserver-express 简版图层发布加载服务
本服务创建的目的为简化三维切片缓存的发布和加载，以更高效的方式完成以上工作。

目前主要用于Cesium（超图改）的测试

功能：
- [ ] 适配超图切片缓存
- [ ] 适配开源切片缓存

## 使用
1. 创建自定义空间，如名称为 testSpace 
2. 空间下创建三维数据，需要指定发布位置，如Config
3. 空间下创建场景，如testScene
4. 场景下创建图层，需要指定使用的三维数据，如 Config
5. Cesium调用

## 特别说明
由于超图数据关系问题，三维数据的名称和图层的名称需要保持一致！(调用接口程序自动处理)

三维数据文件夹中的scp文件名称与三维数据的名称保持一致！

注意，如果使用代理，并且使用了子路径，务必将 /giservices 进行特别处理，因为超图的返回要求是服务端拼接完整的请求路径，程序会将/giservices作为根路径，即 http(s)://yourdomain/giservices/..... 

代理中务必设置 Host的header，否则程序无法获取到真实的请求地址host，如nginx中：
```
proxy_set_header Host $host_name;
proxy_set_header X-Forwarded-Proto $scheme;
```

## 调用地址
适配超图Cesium
scene.open(
'http://localhost:8080/giservices/{自定义空间名称}/rest/realspace'
)

其他数据文件的请求地址：
http://localhost:8080/giservices/{自定义空间名称}/data/{文件}
如
http://localhost:8080/giservices/testSpace/data/roads.kml
