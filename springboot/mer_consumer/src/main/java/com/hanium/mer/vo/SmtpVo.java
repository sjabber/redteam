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
@Entity(name="smtp_info")
public class SmtpVo {

    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)//
    @Column(name = "smtp_no")
    @JsonProperty(value = "smtp_no")
    private Long smtpNo;

    //FK jpa 찾아보기
    @Column(name = "user_no")
    @JsonProperty(value = "user_no")
    private Long userNo;

    @Column(name = "smtp_host")
    @JsonProperty(value = "smtp_host")
    private String smtpHost;

    @Column(name = "smtp_port")
    @JsonProperty(value = "smtp_port")
    private String smtpPort;

    //나중에 이름 변경하기
    @Column(name = "protocol")
    @JsonProperty(value = "smtp_protocol")
    private String smtpProtocol;

    @Column(name = "tls")
    @JsonProperty(value = "smtp_tls")
    private String smtpTls;

    @Column(name = "timeout")
    @JsonProperty(value = "smtp_timeout")
    private String smtpTimeOut;

    @Column(name = "smtp_id")
    @JsonProperty(value = "smtp_id")
    private String smtpId;


    @Column(name = "smtp_pw")
    @JsonProperty(value = "smtp_pw")
    private String smtpPw;

    @Column(name = "created_time")
    @JsonProperty(value = "create_time")
    private LocalDateTime createTime;

    @Column(name = "modified_time")
    @JsonProperty(value = "modify")
    private LocalDateTime modify;

}
