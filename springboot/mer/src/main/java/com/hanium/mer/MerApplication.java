package com.hanium.mer;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

import javax.crypto.Cipher;
import javax.crypto.spec.IvParameterSpec;
import javax.crypto.spec.SecretKeySpec;
import java.security.Key;
import java.util.Base64;

@SpringBootApplication
public class MerApplication {

    public static void main(String[] args) throws Exception {
        SpringApplication.run(MerApplication.class, args);
    }

}
