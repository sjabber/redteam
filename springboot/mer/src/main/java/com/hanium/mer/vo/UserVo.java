package com.hanium.mer.vo;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AccessLevel;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import javax.persistence.*;
import java.time.LocalDateTime;

@Data
@AllArgsConstructor//
@NoArgsConstructor(access = AccessLevel.PROTECTED)//
@Entity(name="user_info")
public class UserVo {


    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)//
    @Column(name = "user_no")
    @JsonProperty(value = "user_no")
    private Long userNo;

    @Column(name = "user_email")
    @JsonProperty(value = "email")
    private String userId;

    @Column(name = "user_name")
    @JsonProperty(value = "name")
    private String userName;

    @Column(name = "user_pw")
    @JsonProperty(value = "user_pw")
    private String userPw;

    @Column(name = "created_time")
    @JsonProperty(value = "created_time")
    private LocalDateTime createdTime;

    @Column(name = "modified_time")
    @JsonProperty(value = "modify_time")
    private LocalDateTime modifyTime;

    @Column(name = "is_enabled")
    @JsonProperty(value = "is_enabled")
    private int isEnabled;
}
