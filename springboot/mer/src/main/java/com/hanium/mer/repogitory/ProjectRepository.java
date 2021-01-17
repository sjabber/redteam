package com.hanium.mer.repogitory;

import com.hanium.mer.vo.ProjectVo;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;


@Repository
public interface ProjectRepository extends JpaRepository<ProjectVo, Long> {

    //TODO targetNo만 있으면됨
    @Query(value = "SELECT target_no as targetNo, target_name as targetName, target_email as targetEmail," +
                    " target_phone as targetPhone, target_organize as targetOrganize, target_position as targetPosition " +
            "FROM target_info " +
            "WHERE user_no = :user_no " +
            "AND (tag1 != 0 AND tag1 IN( :tag1, :tag2, :tag3))" +
            "OR (tag2 != 0 AND tag2 IN( :tag1, :tag2, :tag3))" +
            "OR (tag3 != 0 AND tag3 IN( :tag1, :tag2, :tag3))", nativeQuery = true)
    List<Object[]> findTargetByUserNoAndTags(@Param("user_no") Long userNo,
                                            @Param("tag1") Integer tagFirst,
                                            @Param("tag2") Integer tagSecond,
                                            @Param("tag3") Integer tagThird);

    List<ProjectVo> findByUserNo(Long user_no);

    Optional<ProjectVo> findBypNoAndUserNo(Long p_no, Long user_no);

}
