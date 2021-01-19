package com.hanium.mer.vo;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import javax.persistence.*;
import java.time.LocalDateTime;

@Data
@AllArgsConstructor//
@NoArgsConstructor
@Entity(name="template_info")
public class TemplateVO {

    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)//
    @Column(name = "tmp_no")
    @JsonProperty(value = "tmp_no")
    private Long tmpNo;

    @Column(name = "tmp_division")
    @JsonProperty(value = "tmp_division")
    private int tmpDivision;

    @Column(name = "tmp_kind")
    @JsonProperty(value = "tmp_kind")
    private int tmpKind;

    @Column(name = "file_info")
    @JsonProperty(value = "file_info")
    private int fileInfo;

    @Column(name = "tmp_name")
    @JsonProperty(value = "tmp_name")
    private String tmpName;

    @Column(name = "mail_title")
    @JsonProperty(value = "mail_title")
    private String mailTitle;

    @Column(name = "sender_name")
    @JsonProperty(value = "sender_name")
    private String senderName;

    @Column(name = "download_type")
    @JsonProperty(value = "download_type")
    private int downloadType;

    @Column(name = "created_time")
    @JsonProperty(value = "created_time")
    private LocalDateTime createdTime;

    @Column(name = "modified_time")
    @JsonProperty(value = "modified_time")
    private LocalDateTime modifiedTime;

    @Column(name = "mail_content")
    @JsonProperty(value = "mail_content")
    private String mailContent;

    @Column(name = "user_no")
    @JsonProperty(value = "user_no")
    private Long userNo;
}
