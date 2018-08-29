package tpl

const tplContractPomXML = `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>
    <parent>
        <groupId>mtime.basis</groupId>
        <artifactId>contract-parent</artifactId>
        <version>1.0-SNAPSHOT</version>
        <relativePath></relativePath>
    </parent>
    <groupId>{{.GroupID}}</groupId>
    <artifactId>{{.ArtifactID}}</artifactId>
    <version>1.0-SNAPSHOT</version>             
</project>`

// RPC 服务接口
const tplTestService = `package {{.Package}}.iface;

import {{.Package}}.dto.TestDto.*;

/**
 * 测试服务
**/
public interface TestService {	
	/**
	 * 测试
	**/
	HelloResponse hello(HelloRequest request);
}`

// RPC 服务接口
const tplTestDto = `package {{.Package}}.dto;

import lark.pb.field.FieldType;
import lark.pb.annotation.ProtoField;
import lark.pb.annotation.ProtoMessage;
import {{.Package}}.constant.TestType;

import java.time.LocalDateTime;

public class TestDto {
    /**
     * 请求参数
     */
    @Setter
    @Getter 
    @ProtoMessage(description = "请求参数")
    public static class HelloRequest {
        /**
         * ID
         */
        @ProtoField(order = 1, type = FieldType.INT32, required = true, description = "ID")
        private int id;

        /**
         * 类型
         */
        @ProtoField(order = 2, type = FieldType.ENUM, description = "类型")
        private TestType type;

        /**
         * 时间
         */
        @ProtoField(order = 3, type = FieldType.INT64, description = "时间")
        private LocalDateTime time;
    }

    /**
     * 响应结果
     */
    @Setter
    @Getter     
    @ProtoMessage(description = "响应结果")
    public static class HelloResponse {
        /**
         * 结果
         */
        @ProtoField(order = 1, type = FieldType.STRING, description = "结果")
        private String result;
    }
}
`

// 测试枚举
const tplTestType = `package {{.Package}}.constant;

import lark.core.lang.Description;
import lark.core.lang.EnumValuable;
import lark.core.lang.EnumTitlable;
import lark.core.lang.Enums;

@Description("测试类型")
public enum	TestType implements EnumValuable, EnumTitlable {
    /**
     * 好
     */
    GOOD(1, "好"),
    /**
     * 坏
     */
    BAD(2, "坏");

    private int value;
    private String title;

    private TestType(int value, String title) {
        this.value = value;
        this.title = title;
    }

    /**
     * 获取枚举的 int 值,用于数据保存以及序列化
     *
     * @return 枚举的 int 值
     */
    @Override
    public int value() {
        return this.value;
    }

    /**
     * 获取枚举的显示名称
     *
     * @return 枚举的显示名称
     */
    @Override
    public String title() {
        return this.title;
    }
    
    /**
     * 根据 int 值构建一个枚举对象
     *
     * @param value 需要构建枚举的 int 的值
     * @return 返回相应 value 值的枚举对象
     */
    public static TestType valueOf(int value) {
        return Enums.valueOf(TestType.class, value);
    }
}
`

// RPC 服务描述文件
const tplTestServiceXML = `<?xml version="1.0" encoding="utf-8" ?>
<!--
    javaPackage: 设定生成代码的包路径
    version: 0.2-旧规范, 1.0-新规范(注意 0.2 与 1.0 规范不完全兼容, 不能简单修改版本号)
-->
<rsd javaPackage="com.test.demo.service" version="1.0">
    <!--服务定义
        name: 服务名称(生成代码时会自动加上 Service 后缀)
        alias: 服务别名, 如果指定, 则调用服务时会传递别名到服务端
        description: 服务描述
    -->
    <service name="Test" description="测试服务">
        <!--生成的服务接口需要额外引入的包路径-->
        <imports>
            <!--
                path: 包路径
                lang: 适用的语言, 如 java/c#/go
            -->
            <!--<import path="lark.pb.data.PageList" lang="java"/>-->
        </imports>
        <!--方法定义
            name: 方法名称
            alias: 方法别名, 如果指定, 则调用服务时会传递别名到服务端
            description: 方法描述
        -->
        <method name="Hello" alias="" description="hello">
            <!--请求参数
                multiple: 是否是多参数, 默认 false, 表示单 protobuf 对象
            -->
            <request multiple="false">        
                <field modifier="required" type="int32" name="Id" order="1" description="ID"/>
                <field modifier="optional" type="TestType" name="Type" order="2" kind="enum"/>
                <field modifier="optional" type="int64" name="Time" order="3" javaType="LocalDateTime"/>
            </request>
            <response>
                <field modifier="optional" type="string" name="Result" order="1" description="结果"/>
            </response>
            <errors>
                <!--错误码
                    name: 常量名称
                    code: 常量整数值
                    message: 错误描述
                -->            
                <error name="INVALID_TIME" code="1023" message="时间无效"/>
            </errors>
        </method>
    </service>
    <types>
        <imports>
            <!--<import path="java.math.BigDecimal"/>-->
        </imports>
        <!--类型定义
            name: 类型名称
            description: 类型描述
        -->
        <type name="User" description="用户">
            <!--字段定义
                modifier: 修饰符, 有效值: required/optional/repeated
                type: 字段类型
                name: 字段名称, 各语言应该要自动转换名称的大小写
                order: 顺序号
                kind: 类型补充分类, 有效值: enum/message, 枚举类型必须指定 kind 属性
                description: 字段描述
                javaType: 扩展属性, 表示映射的 Java 类型
                goType: 扩展属性, 表示映射的 Go 类型
            -->
            <field modifier="required" type="int32" name="Id" order="1" description="用户 ID"/>
            <field modifier="required" type="int32" name="Name" order="2" description="用户名称"/>
            <field modifier="required" type="bool" name="Admin" order="3" description="是否是超级管理员"/>
            <field modifier="required" type="UserStatus" name="Status" order="4" description="用户状态" kind="enum"/>
            <field modifier="required" type="int64" name="CreateTime" order="5" javaType="LocalDateTime" goType="time.Time" description="创建时间"/>
            <field modifier="required" type="int64" name="ModifyTime" order="6" javaType="LocalDateTime" goType="time.Time" description="修改时间"/>
        </type>        
    </types>
    <enums>
        <!--枚举定义
            name: 类型名称
            description: 类型描述
        -->    
        <enum name="TestType" description="测试类型">
            <field name="GOOD" value="1" description="好"/>
            <field name="BAD" value="2" description="坏"/>
        </enum>
    </enums>
</rsd>
`

const tplSpringFactories = `# Auto Configure
org.springframework.boot.autoconfigure.EnableAutoConfiguration={{.Package}}.spring.ServiceAutoConfig`

const tplServiceAutoConfig = `package {{.Package}}.spring;

import lark.net.rpc.client.ServiceFactory;
import {{.Package}}.iface.TestService;
import org.springframework.boot.autoconfigure.condition.ConditionalOnMissingBean;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Lazy;
import org.springframework.core.Ordered;
import org.springframework.core.annotation.Order;

@Configuration
@Lazy
@Order(Ordered.LOWEST_PRECEDENCE)
public class ServiceAutoConfig {
    private final String SERVER = "{{.Package}}";

    @Bean
    @ConditionalOnMissingBean
    public TestService testService(ServiceFactory serviceFactory) {
        return serviceFactory.get(SERVER, TestService.class);
    }
}
`
