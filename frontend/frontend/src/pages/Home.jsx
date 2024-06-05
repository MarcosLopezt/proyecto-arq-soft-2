import React, { useState, useEffect } from "react";
import "./Home.css";
import Courses from "../components/Courses";
import { useLocation } from "react-router-dom";
import {
  AppBar,
  Toolbar,
  Typography,
  IconButton,
  InputBase,
  Button,
  Avatar,
  Menu,
  List,
  ListItemIcon,
  ListItemText,
  ListItemButton,
} from "@mui/material";
import SearchIcon from "@mui/icons-material/Search";
import AccountCircleIcon from "@mui/icons-material/AccountCircle";
import LogoutIcon from "@mui/icons-material/Logout";
import { useNavigate } from "react-router-dom";
import "../components/Componentes.css";
import { useFormik } from "formik";
import * as Yup from "yup";
//import { SearchBar } from './SearchBar';
//validar permisos para ver que modulos mostramos en la navbar

function Home() {
  const location = useLocation();
  const role = location.state?.role;

  const navigate = useNavigate();
  const [logoutOpen, setLogoutOpen] = useState(false);
  const [userRole, setUserRole] = useState("");
  //const [searchText, setSearchText] = useState("");
  const [courses, setCourses] = useState([]);

  useEffect(() => {
    setUserRole(role);
  }, [role]);

  const handleLogoutClick = () => {
    setLogoutOpen(true);
  };

  const handleLogoutClose = () => {
    setLogoutOpen(false);
  };

  const logout = () => {
    document.cookie =
      "session_token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT;";
    navigate("/");
  };

  // const handleChange = (event) => {
  //   setSearchText(event.target.value);
  // };

  const submit = (values) => {
    search(values.text);
  };

  const { handleSubmit, handleChange, values, errors } = useFormik({
    initialValues: {
      text: "",
    },

    validationSchema: Yup.object({
      text: Yup.string()
        .required("¡Campo Requerido!")
        .max(255, "Maximo 255 caracteres"),
    }),

    onSubmit: submit,
  });

  const search = async (name) => {
    //console.log(name);
    const response = await fetch(`http://localhost:8080/cursos/${name}`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
      //body: JSON.stringify({ "email": email, "password": password })
    });

    if (response.status === 200) {
      const data = await response.json();
      //console.log(data[0]);
      setCourses(data);
      console.log(courses[0]);
    } else {
      console.log("No existe el curso");
    }
  };

  const handleButtonClick = () => {
    navigate("/createCourse");
  };

  return (
    <>
      <AppBar
        className="navbar"
        position="static"
        sx={{ backgroundColor: "#785589" }}
      >
        <Toolbar>
          {/* Logo */}
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            Logo
          </Typography>

          {/* Barra de búsqueda */}
          <div style={{ marginRight: "20px" }}>
            <div
              style={{
                display: "flex",
                alignItems: "center",
                backgroundColor: "white",
                borderRadius: "4px",
                paddingLeft: "10px",
              }}
            >
              <SearchIcon />
              <form onSubmit={handleSubmit}>
                <InputBase
                  placeholder="Buscar..."
                  inputProps={{ "aria-label": "buscar" }}
                  style={{ marginLeft: "10px" }}
                  name="text"
                  value={values.text}
                  onChange={handleChange}
                  error={!!errors.text}
                  helperText={errors.text}
                />
              </form>
            </div>
          </div>

          {/* Botón de "Mis Cursos" */}
          <Button className="button-misCursos" variant="contained">
            Mis Cursos
          </Button>

          {/* Botón de "Crear Curso" */}
          {userRole === "admin" && (
            <Button
              className="button-crear-curso"
              variant="contained"
              onClick={handleButtonClick}
            >
              Crear Curso
            </Button>
          )}

          {/* Icono de perfil */}
          <IconButton id="profile-icon" onClick={handleLogoutClick}>
            <Avatar>
              <AccountCircleIcon />
            </Avatar>
          </IconButton>

          {/* Menú de logout */}
          <Menu
            anchorEl={
              logoutOpen ? document.getElementById("profile-icon") : null
            }
            open={logoutOpen}
            onClose={handleLogoutClose}
          >
            <div className="list-conteiner">
              <List>
                <ListItemButton onClick={logout} className="button-logout">
                  <ListItemIcon>
                    <LogoutIcon className="icon-logout" />
                  </ListItemIcon>
                  <ListItemText primary="Logout" />
                </ListItemButton>
              </List>
            </div>
          </Menu>
        </Toolbar>
      </AppBar>
      <Courses courses={courses}></Courses>
    </>
  );
}

export default Home;
