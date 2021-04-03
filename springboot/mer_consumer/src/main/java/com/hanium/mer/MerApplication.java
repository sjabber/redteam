package com.hanium.mer;

import lombok.extern.slf4j.Slf4j;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.WebApplicationType;
import org.springframework.boot.autoconfigure.SpringBootApplication;

import java.io.File;

@Slf4j
@SpringBootApplication
public class MerApplication {

    public static void main(String[] args) {
        String path = "C:\\mailTemp\\"; //폴더 경로
        File Folder = new File(path);

        // 해당 디렉토리가 없을경우 디렉토리를 생성합니다.
        if (!Folder.exists()) {
            try{
                Folder.mkdir(); //폴더 생성합니다.
                log.info("make mailTemp");
            }
            catch(Exception e){
                e.getStackTrace();
            }
        }else {

        }

        SpringApplication springApplication = new SpringApplication(MerApplication.class);
        springApplication.setWebApplicationType(WebApplicationType.NONE);
        springApplication.run(args);
    }

}
