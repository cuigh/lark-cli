package {{.Package}};

import {{.Package}}.constant.TestType;
import {{.Package}}.dto.TestDto;
import {{.Package}}.iface.TestService;
import org.junit.Test;

import java.time.LocalDateTime;

public class TestServiceTests {
    @Autowired
    private TestService testService;

    @Test
    public void testHello() throws Exception {
        TestDto.HelloRequest request = new TestDto.HelloRequest();
        request.setId(123);
        request.setType(TestType.GOOD);
        request.setTime(LocalDateTime.now());

        TestDto.HelloResponse response = testService.hello(request);
        printJsonln(response);
    }
}