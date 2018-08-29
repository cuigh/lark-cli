package tpl

// 项目 pom 文件模板
const tplProjectPomXML = `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>
    <groupId>{{.GroupID}}</groupId>
    <artifactId>{{.ArtifactID}}</artifactId>
    <packaging>pom</packaging>
    <version>1.0-SNAPSHOT</version>
    <modules>
    </modules>

    <build>
        <plugins>
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

// 自述文件模板
const tplReadMe = `# {{.ArtifactID}}

[项目描述...]`

// Git 忽略文件模板
const tplGitIgnore = `*.class

# Mobile Tools for Java (J2ME)
.mtj.tmp/

# Package Files #
*.jar
*.war
*.ear

# virtual machine crash logs
hs_err_pid*

# Build results
target

# IDEA files
.idea
*.iml

# Eclipse files
.settings
.classpath
.project
`
