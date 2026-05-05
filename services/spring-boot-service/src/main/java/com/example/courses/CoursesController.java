package com.example.courses;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import java.util.Collections;

@RestController
@RequestMapping("/api/v1")
public class CoursesController {

    @GetMapping("/courses")
    public java.util.List<Object> getCourses() {
        return Collections.emptyList();
    }
}
