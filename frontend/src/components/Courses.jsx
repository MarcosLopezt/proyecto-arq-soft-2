import React from "react";
import CourseCard from "./CourseCard";
import "./CourseCard.css";

function capitalizeTitle(title) {
  return title.replace(/\b\w/g, (char) => char.toUpperCase());
}

function Courses({ courses }) {
  //console.log("Courses:", courses);
  console.log("Courses recibió:", courses);
  console.log("Primer curso:", courses[0]);
  
  const modulos = courses.length;
  const shouldWrap = modulos > 4;
  //const margenEntreTarjetas = `calc(100% / ${modulos})`;
  if (modulos === 0) {
    return null;
  }

  return (
    <div
      className={`contenedor-cards ${shouldWrap ? "wrap" : ""}`}
      style={{ padding: "60px" }}
    >
      {courses.map((course) => {
        console.log("Mapeando curso:", course);
        console.log("course.ID:", course.ID);
        console.log("course.id:", course.id);
        console.log("course.course_id:", course.course_id);
        
        return (
          <CourseCard
            key={course.ID || course.id || course.course_id}
            ID={course.ID || course.id || course.course_id}
            title={capitalizeTitle(course.course_name)}
            description={course.description}
            category={course.category}
            length={course.length}
            cupos={course.cupos}
            modulos={modulos}
            image={`/assets/skillup${Math.floor(Math.random() * 6) + 2}.jpg`}
          />
        );
      })}
    </div>
  );
}

export default Courses;
