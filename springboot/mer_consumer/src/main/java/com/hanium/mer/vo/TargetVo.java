package com.hanium.mer.vo;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import javax.persistence.*;
import java.time.LocalDateTime;

@Data
@AllArgsConstructor//
@NoArgsConstructor//
@Entity(name="target_info")
public class TargetVo {

    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)//
    @Column(name = "target_no")
    @JsonProperty(value = "target_no")
    private int targetNo;

    @Column(name = "target_name")
    @JsonProperty(value = "target_name")
    private String targetName;

    @Column(name = "target_email")
    @JsonProperty(value = "target_email")
    private String targetEmail;

    @Column(name = "target_phone")
    @JsonProperty(value = "target_phone")
    private String targetPhone;

    @Column(name = "target_organize")
    @JsonProperty(value = "target_organize")
    private String targetOrganize;

    @Column(name = "target_position")
    @JsonProperty(value = "target_position")
    private String targetPosition;

    @Column(name = "created_time")
    @JsonProperty(value = "created_time")
    private LocalDateTime createdTime;

    @Column(name = "modified_time")
    @JsonProperty(value = "modified_time")
    private LocalDateTime modifiedTime;

    @Column(name = "user_no")
    @JsonProperty(value = "user_no")
    private Long userNo;

    @Column(name = "tag1")
    @JsonProperty(value = "tag1")
    private int tagFirst;

    @Column(name = "tag2")
    @JsonProperty(value = "tag2")
    private int tagSecond;

    @Column(name = "tag3")
    @JsonProperty(value = "tag3")
    private int tagThird;

}
