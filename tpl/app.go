package tpl

// service/msg/task/web 模块 pom.xml 文件模板
const tplAppPomXML = `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>
    <parent>
        <groupId>mtime.basis</groupId>
        <artifactId>{{.Type}}-parent</artifactId>
        <version>1.0-SNAPSHOT</version>
        <relativePath></relativePath>
    </parent>
    <groupId>{{.GroupID}}</groupId>
    <artifactId>{{.ArtifactID}}</artifactId>
    <version>1.0-SNAPSHOT</version>
    {{- if eq .Type "service"}}
    
    <dependencies>
        <dependency>
            <groupId>{{.GroupID}}</groupId>
            <artifactId>{{.ArtifactID}}-contract</artifactId>
            <version>1.0-SNAPSHOT</version>
        </dependency>
    </dependencies>
    {{end -}}
    
    <build>
        <plugins>
            <plugin>
                <groupId>org.apache.maven.plugins</groupId>
                <artifactId>maven-assembly-plugin</artifactId>
            </plugin>
            <plugin>
                <groupId>org.apache.maven.plugins</groupId>
                <artifactId>maven-deploy-plugin</artifactId>
                <configuration>
                    <skip>true</skip>
                </configuration>
            </plugin>
        </plugins>
    </build>

</project>`

// 构建信息文件模板(build.json)
const tplBuildJSON = `{
  "productVersion": "${project.version}",
  "scmRevision": "${git.commit.id}",
  "scmBranch": "${git.branch}",
  "scmUser": "${git.commit.user.name}",
  "scmMessage": "${git.commit.message.short}",
  "scmTime": "${git.commit.time}",
  "buildAt": "${git.build.time}",
  "buildBy": "${user.name}",
  "buildOn": "${git.build.host}"
}`

// 打包配置文件模板（assembly.xml）
const tplAssemblyXML = `<assembly xmlns="http://maven.apache.org/plugins/maven-assembly-plugin/assembly/1.1.3"
          xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
          xsi:schemaLocation="http://maven.apache.org/plugins/maven-assembly-plugin/assembly/1.1.3 http://maven.apache.org/xsd/assembly-1.1.3.xsd">
    <id>pkg</id>
    <formats>
        <format>dir</format>
    </formats>
    <files>
        <file>
            <source>build/build.json</source>
            <outputDirectory></outputDirectory>
            <filtered>true</filtered>
        </file>
        <file>
            <source>build/pom.properties</source>
            <outputDirectory></outputDirectory>
            <filtered>true</filtered>
        </file>
    </files>
    <fileSets>
        <fileSet>
            <directory>src/main/resources</directory>
            <outputDirectory></outputDirectory>
        </fileSet>
    </fileSets>
    <dependencySets>
        <dependencySet>
            <useProjectArtifact>true</useProjectArtifact>
            <scope>runtime</scope>
            <outputDirectory>lib</outputDirectory>
        </dependencySet>
    </dependencySets>
</assembly>`

// 程序信息文件模板(pom.properties)
const tplPomProperties = `version=${project.version}
groupId=${project.groupId}
artifactId=${project.artifactId}
mainClass={{.Package}}.Bootstrap`

// app.conf 配置文件
const tplAppConfig = `<config>
    <app>
        <add key="app_name" value="{{.Package}}"/>
        <add key="debug" value="true"/>
        <add key="monitor_enabled" value="true"/>
        <add key="monitor_port" value="[请替换成实际的监控端口]"/>
    </app>
</config>
`

// 业务类模板
const tplTestBiz = `package {{.Package}}.biz;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import {{.Package}}.dao.TestDao;
import {{.Package}}.entity.TestObject;

@Service
public class TestBiz {
	@Autowired
	private TestDao testDao;

    // todo: remove this method
	public TestObject getObject(int id) {
        return testDao.getObject(id);
	}
}
`

// 数据访问类模板
const tplTestDao = `package {{.Package}}.dao;

import {{.Package}}.entity.TestObject;
import org.springframework.stereotype.Repository;

@Repository
public class TestDao {
    // todo: remove this method
    public TestObject getObject(int id) {
        TestObject object = new TestObject();
        object.setId(id);
        object.setName("noname");
        return object;
    }
}
`

// 实体类模板
const tplTestEntity = `package {{.Package}}.entity;

import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class TestObject {
    private int id;
    private String name;
}
`

// 抽象测试基类模板
const tplAbstractTest = `package {{.Package}};

import lark.util.test.TestBase;
import lark.util.test.spring.SpringJUnit;
import org.junit.BeforeClass;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.autoconfigure.data.mongo.MongoDataAutoConfiguration;
import org.springframework.boot.autoconfigure.mongo.MongoAutoConfiguration;

public abstract class AbstractTest extends TestBase {
    @BeforeClass
    public static void init() {
        SpringJUnit.boot(Dummy.class, Bootstrap.class);
    }

    @SpringBootApplication
    public static class Dummy {
    }
}
`
