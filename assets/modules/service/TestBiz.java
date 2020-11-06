package {{.Package}}.biz;

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