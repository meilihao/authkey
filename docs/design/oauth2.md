# 理解 OAuth 2.0 认证流程
OAuth 2.0 的标准是 RFC 6749 文件, OAuth 的核心就是向第三方应用颁发令牌.

OAuth 引入了一个授权层，用来分离两种不同的角色：客户端和资源所有者. 经资源所有者同意以后，资源服务器可以向客户端颁发令牌. 客户端通过令牌，去请求数据.

OAuth 2.0 共有 4 种访问模式：
- 授权码模式(Authorization Code)，这种方式是最常用的流程，安全性也最高，它适用于那些有后端的 Web 应用
- 简化模式(Implicit)，适用于纯网页端应用，不再使用

    Implicit Flow 最大的弊端就是无法使用刷新令牌（Refresh Token），每次访问令牌过期都需要重新发起授权流程去获取新的访问令牌. 对于 Implicit Flow 原本针对的场景，最新 OAuth2 标准推荐使用标准的 Authorization Code 结合 [PKCE](https://auth0.com/docs/flows/concepts/auth-code-pkce) 扩展来实现认证授权
- 密码模式(Resource owner password credentials), 不介绍
- 客户端模式(Client credentials), 适用于没有前端的命令行应用, 不介绍

## OAuth2授权模式的选型
![按授权需要的多端情况](docs/misc/img/1149398-20191203120032871-884830528.png)
![按客户端类型与所有者](docs/misc/img/1149398-20191203140907868-924472362.png)

## Authorization Code
Authorization Code具体流程：
1. 用户打开客户端以后，客户端要求用户给予授权
1. 用户同意给予客户端授权
1. 客户端使用上一步获得的授权，向认证服务器申请令牌
1. 认证服务器对客户端进行认证以后，确认无误，同意发放令牌
1. 客户端使用令牌，向资源服务器申请获取资源
1. 资源服务器确认令牌无误，同意向客户端开放资源

![](docs/misc/img/oauth-authorization-code.svg)