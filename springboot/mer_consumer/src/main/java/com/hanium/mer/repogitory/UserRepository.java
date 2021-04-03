package com.hanium.mer.repogitory;

import com.hanium.mer.vo.UserVo;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;

@Repository
public interface UserRepository extends JpaRepository<UserVo, Long> {

    public List<UserVo> findAll();

    public Optional<UserVo> findByUserNo(Long user_no);

}
