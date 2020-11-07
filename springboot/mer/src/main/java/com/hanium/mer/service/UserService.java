package com.hanium.mer.service;

import com.hanium.mer.repogitory.UserRepository;
import com.hanium.mer.vo.NewUserInfo;
import com.hanium.mer.vo.UserVo;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;

@Service
public class UserService {

    @Autowired
    private UserRepository userRepository;

    public List<UserVo> findAll(){
        List<UserVo> users = new ArrayList<>();
        userRepository.findAll().forEach(e -> users.add(e));
        return users;
    }

    public Optional<UserVo> findByUserNo(Long user_no){
        Optional<UserVo> user =  userRepository.findByUserNo(user_no);
        return user;
    }

    public void changeUserInfo(Long user_no, NewUserInfo newUserInfo){
        UserVo user =  userRepository.findByUserNo(user_no).get();
        user.setUserId(newUserInfo.getEmail());
        user.setUserName(newUserInfo.getName());
        user.setUserPw(newUserInfo.getPassword());
        user.setModifyTime(LocalDateTime.now());
        System.out.println(userRepository.save(user));
        //error 는?
    }

}
