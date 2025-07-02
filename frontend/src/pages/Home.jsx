import React, { useState, useEffect } from "react";
import "./Home.css";
import Courses from "../components/Courses";
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
  Snackbar,
  Alert,
  FormControlLabel,
  Checkbox,
  Box,
} from "@mui/material";
import SearchIcon from "@mui/icons-material/Search";
import SettingsIcon from "@mui/icons-material/Settings";
import AccountCircleIcon from "@mui/icons-material/AccountCircle";
import LogoutIcon from "@mui/icons-material/Logout";
import { useNavigate } from "react-router-dom";
import "../components/Componentes.css";
import { useFormik } from "formik";
import * as Yup from "yup";
import { isAuthenticated } from "../utils/authUtils";

function Home() {
  //const location = useLocation();
  const userRole = localStorage.getItem("userRole");

  const navigate = useNavigate();
  const [logoutOpen, setLogoutOpen] = useState(false);
  //const [userRole, setUserRole] = useState(role);
  const [courses, setCourses] = useState([]);
  const [recomendados, setRecomendados] = useState(true);
  const [authenticated, setAuthenticated] = useState(false);
  const [snackbarOpen, setSnackbarOpen] = useState(false);
  const [snackbarMessage, setSnackbarMessage] = useState("");
  const [availableOnly, setAvailableOnly] = useState(false);

  // useEffect(() => {
  //   setUserRole(role);
  // }, [role]);

  useEffect(() => {
    if (recomendados) {
      searchRecommended();
    }
  }, [recomendados, availableOnly]);

  useEffect(() => {
    setAuthenticated(isAuthenticated);
    if (authenticated) {
      logout();
    }
  }, []);

  const handleLogoutClick = () => {
    setLogoutOpen(true);
  };

  const handleLogoutClose = () => {
    setLogoutOpen(false);
  };

  const handleSnackbarClose = () => {
    setSnackbarOpen(false);
  };

  const logout = () => {
    localStorage.removeItem("authToken");
    navigate("/");
  };

  const searchRecommended = async () => {
    // Siempre usar el servicio de búsqueda, con availableOnly según el estado del checkbox
    const availableParam = availableOnly ? "&availableOnly=true" : "&availableOnly=false";
    const response = await fetch(
      `http://localhost:8085/search?q=*:*&offset=0&limit=10${availableParam}`,
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      }
    );

    if (response.status === 200) {
      const data = await response.json();
      console.log("Respuesta completa del servicio de búsqueda:", data);
      console.log("data.courses:", data.courses);
      
      if (!data.courses || data.courses.length === 0) {
        setSnackbarMessage("No existen cursos disponibles.");
        setSnackbarOpen(true);
      } else {
        console.log("Primer curso de la respuesta:", data.courses[0]);
        setCourses(data.courses);
      }
    } else {
      console.log("Error al buscar cursos");
      setSnackbarMessage("Error al buscar cursos.");
      setSnackbarOpen(true);
    }
  };

  const submit = (values) => {
    if (!values.text.trim()) {
      setRecomendados(true);
      searchRecommended(); //muestra cursos recomendado (por defecto de informatica)
    } else {
      setRecomendados(false);
      search(values.text);
    }
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
    const availableParam = availableOnly ? "&availableOnly=true" : "";
    const response = await fetch(
      `http://localhost:8085/search?q=${name}&offset=0&limit=10${availableParam}`,
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      }
    );

    if (response.status === 200) {
      const data = await response.json();
      console.log("Respuesta completa de búsqueda por nombre:", data);
      console.log("data.courses:", data.courses);
      
      if (!data.courses || data.courses.length === 0) {
        setSnackbarMessage("No existen cursos con ese nombre o categoría.");
        setSnackbarOpen(true);
      } else {
        console.log("Primer curso de la búsqueda:", data.courses[0]);
        setCourses(data.courses);
      }
    } else {
      console.log("no encontre cursos");
      setSnackbarMessage("No existen cursos con ese nombre o categoría.");
      setSnackbarOpen(true);
    }
  };

  const handleButtonClick = () => {
    setRecomendados(true);
    navigate("/createCourse");
  };

  const handleMisCursosButton = () => {
    setRecomendados(true);
    navigate("/mycourses");
  };

  const navigateToMicroservicios = () => {
    navigate("/microservicios"); // Navegar a la ruta '/microservicios'
  };

  const handleAvailableOnlyChange = (event) => {
    setAvailableOnly(event.target.checked);
    // Si no hay texto de búsqueda, ejecutar búsqueda automáticamente
    if (!values.text.trim()) {
      setRecomendados(true);
    }
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
          <Button
            className="button-misCursos"
            variant="contained"
            onClick={handleMisCursosButton}
          >
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
              {/* Botón para navegar a /microservicios con icono de ajustes */}
              {userRole === "admin" && (
              <List>
                <ListItemButton onClick={navigateToMicroservicios}>
                  <ListItemIcon>
                    <SettingsIcon /> {/* Icono de ajustes */}
                  </ListItemIcon>
                  <ListItemText primary="Ir a Microservicios" />
                </ListItemButton>
              </List>
              )}
            </div>
          </Menu>
        </Toolbar>
      </AppBar>

      {/* Checkbox para filtrar por cupos disponibles */}
      <Box sx={{ padding: "20px 60px 0px 60px" }}>
        <FormControlLabel
          control={
            <Checkbox
              checked={availableOnly}
              onChange={handleAvailableOnlyChange}
              sx={{
                color: "#785589",
                "&.Mui-checked": {
                  color: "#785589",
                },
              }}
            />
          }
          label="Mostrar solo cursos con cupos disponibles"
          sx={{
            color: "#785589",
            fontWeight: "bold",
          }}
        />
      </Box>

      {courses && courses.length > 0 && <Courses courses={courses} />}

      <Snackbar
        open={snackbarOpen}
        autoHideDuration={6000}
        onClose={handleSnackbarClose}
        anchorOrigin={{ vertical: "top", horizontal: "right" }}
      >
        <Alert
          onClose={handleSnackbarClose}
          severity="error"
          sx={{
            width: "100%",
            fontSize: "1.2em",
            padding: "20px",
            maxWidth: "600px",
          }}
        >
          {snackbarMessage}
        </Alert>
      </Snackbar>
    </>
  );
}

export default Home;
