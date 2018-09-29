package {{.Package}}.controller;

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