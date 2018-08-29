package tpl

// Web 程序启动类模板
const tplWebBootstrap = `package {{.Package}};

import lark.web.WebApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.autoconfigure.data.mongo.MongoDataAutoConfiguration;
import org.springframework.boot.autoconfigure.mongo.MongoAutoConfiguration;

@SpringBootApplication(exclude = {MongoAutoConfiguration.class, MongoDataAutoConfiguration.class})
public class Bootstrap {
    public static void main(String[] args) {
        WebApplication app = new WebApplication(Bootstrap.class, args);
        app.run();
    }
}
`

// Web 程序控制器类模板
const tplTestController = `package {{.Package}}.controller;

import {{.Package}}.biz.TestBiz;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.*;

/**
 * 测试
 */
@Controller
@RequestMapping("/")
public class TestController {
    @Autowired
    private TestBiz testBiz;
    
    @RequestMapping(value = "hello", method = RequestMethod.GET)
    @ResponseBody
    public String hello(){
        return "Hello, world.";
    }
}
`

// app.conf 配置文件
const tplWebAppConfig = `<config>
    <app>
        <add key="app_name" value="{{.Package}}"/>
        <add key="debug" value="true"/>
    </app>
    <web>
        <add key="auth.form.enabled" value="false"/>
    </web>
</config>
`

// Web 程序控制器类模板
const tplWebAppProperties = `# 监听端口
server.port=8080
`
