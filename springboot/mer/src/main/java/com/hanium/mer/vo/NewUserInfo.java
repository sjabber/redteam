package com.hanium.mer.vo;

import lombok.AccessLevel;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;


@Data
@AllArgsConstructor//
@NoArgsConstructor(access = AccessLevel.PROTECTED)//
public class NewUserInfo {

    private String email;

    private String name;

    private String current_pw;

    private String password;

    private String password_check;
}
