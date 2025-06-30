import React, { useState, useEffect } from "react";
import "./Home.css";
import {
  AppBar,
  Toolbar,
  Typography,
  IconButton,
  Button,
  Avatar,
  Menu,
  List,
  ListItemIcon,
  ListItemText,
  ListItemButton,
  Container,
  Grid,
  Paper,
  Snackbar,
  Alert,
  InputBase,
  Rating,
} from "@mui/material";
import AccountCircleIcon from "@mui/icons-material/AccountCircle";
import LogoutIcon from "@mui/icons-material/Logout";
import { useNavigate } from "react-router-dom";
import Comments from "./Coments";
import Files from "./Files";

import "../components/Componentes.css";

function Course() {
  const [logoutOpen, setLogoutOpen] = useState(false);
  const navigate = useNavigate();
  
  // Agregar logs para debug del courseID
  const courseIDRaw = localStorage.getItem("courseID");
  console.log("courseID raw del localStorage:", courseIDRaw);
  console.log("tipo de courseID raw:", typeof courseIDRaw);
  
  const courseID = parseInt(courseIDRaw, 10);
  console.log("courseID después de parseInt:", courseID);
  console.log("tipo de courseID después de parseInt:", typeof courseID);
  
  const titulo = localStorage.getItem("cursoTitulo");
  const descripcion = localStorage.getItem("cursoDescripcion");
  const categoria = localStorage.getItem("cursoCategoria");
  const length = localStorage.getItem("cursoLength");
  const [open, setOpen] = useState(false);
  const [subscripto, setSubscripto] = useState(false);
  const [snackSubscribed, setSnackSubscribed] = useState(false);
  const userRole = localStorage.getItem("userRole");
  const [comentarioText, setComentarioText] = useState("");
  const userID = localStorage.getItem("userID");
  const [snackComent, setSnackComent] = useState(false);
  const [comentErr, setComentErr] = useState(false);
  const [value, setValue] = useState(0);
  const [disp, setDisp] = useState(0);
  const [disponibles, setDisponibles] = useState(true);
  const [loadingAvailability, setLoadingAvailability] = useState(false);
  const [availabilityError, setAvailabilityError] = useState(null);

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

  useEffect(() => {
    console.log("useEffect ejecutado - courseID:", courseID);
    if (courseID) {
      console.log("Ejecutando searchDisp con courseID:", courseID);
      searchDisp();
    } else {
      console.log("No hay courseID disponible");
    }
  }, [courseID]);

  const searchDisp = async () => {
    console.log("searchDisp iniciado - courseID:", courseID);
    if (!courseID) {
      console.log("No hay courseID, saliendo de searchDisp");
      return;
    }
    
    console.log("Iniciando petición de disponibilidad...");
    setLoadingAvailability(true);
    setAvailabilityError(null);
    
    try {
      const response = await fetch(
        `http://localhost:8083/cursos/availability/concurrent`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify([courseID]), // Enviar array con el ID del curso
        }
      );

      if (response.ok) {
        const data = await response.json();
        console.log("Respuesta disponibilidad concurrente:", data);
        
        // Buscar el resultado para el curso actual
        const courseResult = data.results.find(result => result.course_id === courseID);
        
        if (courseResult) {
          setDisp(courseResult.disponibles || 0);
          console.log(`Curso ${courseID} - Disponibles: ${courseResult.disponibles}`);
          console.log(`Información completa del curso:`, courseResult);
        } else {
          setDisp(0);
          setAvailabilityError("No se pudo obtener la disponibilidad del curso");
        }
      } else {
        console.error("Error al obtener disponibilidad:", response.status);
        setAvailabilityError("Error al calcular la disponibilidad");
        setDisp(0);
      }
    } catch (error) {
      console.error("Error en la petición de disponibilidad:", error);
      setAvailabilityError("Error de conexión al calcular disponibilidad");
      setDisp(0);
    } finally {
      setLoadingAvailability(false);
    }
  };

  const handleSubscription = async () => {
    const userID = parseInt(localStorage.getItem("userID"), 10);
    //console.log(userID);
    //console.log(courseID);

    const response = await fetch(`http://localhost:8084/subscriptions/sub`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ user_id: userID, course_id: courseID }),
    });

    if (response.status === 201) {
      const data = await response.json();
      console.log(data);
      setOpen(true);
      setSubscripto(true);
      
      // Actualizar la disponibilidad después de la suscripción exitosa
      setTimeout(() => {
        searchDisp();
      }, 1000); // Pequeño delay para asegurar que la suscripción se procese
    } else {
      console.log("Error en la subscripcion");
    }
  };

  const handleDeleteButton = async () => {
    const response = await fetch(`http://localhost:8083/cursos/delete`, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ ID: courseID }),
    });

    console.log(response);

    if (response.status === 200) {
      const data = await response.json();
      console.log(data);
      navigate("/home");
    } else {
      console.log("Error al borrar");
    }
  };

  useEffect(() => {
    if (!subscripto) {
      isSubscribed();
    }
  }, [subscripto]);

  const isSubscribed = async () => {
    const userID = parseInt(localStorage.getItem("userID"), 10);
    const response = await fetch(
      `http://localhost:8084/subscriptions/get/${userID}`,
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      }
    );

    if (response.status === 200) {
      const cursos = await response.json();
      const cursoEncontrado = cursos.find(
        (curso) => courseID === curso.course_id
      );
      if (cursoEncontrado) {
        setSubscripto(true);
      }
    } else {
      console.log("No esta suscripto a ningun curso");
    }
  };

  const handleComment = async () => {
    const userintID = parseInt(userID);
    const valueInt = parseInt(value);

    if (comentarioText === "") {
      setComentErr(true);
      return;
    }

    //console.log(valueInt);
    const response = await fetch(`http://localhost:8084/coments/coment`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        user_id: userintID,
        curso_id: courseID,
        texto: comentarioText,
        valor: valueInt,
      }),
    });

    console.log(response);

    if (response.status === 201) {
      //const data = await response.json();
      setSnackComent(true);
    } else {
      console.log("Error al comentar");
    }
  };

  const handleClose = (reason) => {
    if (reason === "clickaway") {
      return;
    }

    setOpen(false);
    setSnackSubscribed(false);
    setSnackComent(false);
    setComentErr(false);
  };

  const handleMisCursosButton = () => {
    navigate("/mycourses");
  };

  const handleSubscribed = () => {
    setSnackSubscribed(true);
  };

  const handleUpdateButton = () => {
    navigate("/update");
  };

  const handleAddFile = () => {
    navigate("/upload");
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
          <Button
            onClick={() => navigate("/home")}
            className="logo-button"
            sx={{
              color: "inherit",
              textTransform: "none",
            }}
          >
            <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
              Logo
            </Typography>
          </Button>

          <div style={{ flexGrow: 1 }}></div>

          {/* Botón de "Mis Cursos" */}
          <Button
            className="button-misCursos"
            variant="contained"
            onClick={handleMisCursosButton}
          >
            Mis Cursos
          </Button>

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
      <Container maxWidth="lg" sx={{ marginTop: "70px" }}>
        <Paper
          sx={{
            padding: "20px",
            border: "3px solid #785589",
            backgroundColor: "#f0f0f0",
          }}
        >
          <Typography variant="h4" gutterBottom>
            {titulo}
          </Typography>
          <Grid container spacing={2}>
            <Grid item xs={12} md={6}>
              <Typography variant="body1">Categoría: {categoria}</Typography>
            </Grid>
            <Grid item xs={12} md={6}>
              <Typography variant="body1">
                Duracion: {length} semanas
              </Typography>
            </Grid>
          </Grid>
          <Typography variant="body1" sx={{ marginTop: "20px" }}>
            Descripción:
          </Typography>
          <Typography
            variant="body2"
            sx={{ whiteSpace: "pre-line", marginTop: "10px" }}
          >
            {descripcion}
          </Typography>
          
          {/* Sección de disponibilidad mejorada */}
          <Typography variant="h6" sx={{ marginTop: "20px", marginBottom: "10px" }}>
            Disponibilidad:
          </Typography>
          
          {loadingAvailability ? (
            <Typography variant="body1" sx={{ color: "#666", fontStyle: "italic" }}>
              Calculando disponibilidad...
            </Typography>
          ) : availabilityError ? (
            <Typography variant="body1" sx={{ color: "#d32f2f", marginBottom: "10px" }}>
              {availabilityError}
            </Typography>
          ) : (
            <Typography variant="body1" sx={{ 
              color: disp > 0 ? "#2e7d32" : "#d32f2f",
              fontWeight: "bold",
              marginBottom: "10px"
            }}>
              {disp > 0 ? `${disp} cupo(s) disponible(s)` : "No hay cupos disponibles"}
            </Typography>
          )}
          
          
          <Button
            variant="contained"
            className="button-subscribe"
            onClick={subscripto ? handleSubscribed : handleSubscription}
            disabled={disp === 0 || loadingAvailability}
            sx={{
              marginTop: "20px",
              backgroundColor: disp === 0 || loadingAvailability ? "gray" : "rgb(49, 45, 45)",
            }}
          >
            {loadingAvailability 
              ? "Calculando disponibilidad..." 
              : disp === 0
              ? "No hay disponibilidad"
              : subscripto
              ? "Inscripto"
              : "Inscribirme ahora"}
          </Button>

          {userRole === "admin" && (
            <>
              <Button
                className="button-add-file"
                variant="contained"
                onClick={handleAddFile}
              >
                Agregar Archivos
              </Button>
              <Button
                className="button-editar-curso"
                variant="contained"
                onClick={handleUpdateButton}
              >
                Editar Curso
              </Button>

              <Button
                className="button-delete-curso"
                variant="contained"
                onClick={handleDeleteButton}
              >
                Eliminar Curso
              </Button>
            </>
          )}
          <Snackbar
            open={open}
            autoHideDuration={6000}
            onClose={handleClose}
            anchorOrigin={{ vertical: "top", horizontal: "right" }}
          >
            <Alert
              onClose={handleClose}
              severity="success"
              sx={{
                width: "100%",
                fontSize: "1.2em",
                padding: "20px",
                maxWidth: "600px",
              }}
            >
              Subscripción realizada con éxito
            </Alert>
          </Snackbar>
          <Snackbar
            open={snackSubscribed}
            autoHideDuration={6000}
            onClose={handleClose}
            anchorOrigin={{ vertical: "top", horizontal: "right" }}
          >
            <Alert
              onClose={handleClose}
              severity="error"
              sx={{
                width: "100%",
                fontSize: "1.2em",
                padding: "20px",
                maxWidth: "600px",
              }}
            >
              ¡Ya estás subscripto!
            </Alert>
          </Snackbar>
          <Snackbar
            open={snackComent}
            autoHideDuration={6000}
            onClose={handleClose}
            anchorOrigin={{ vertical: "top", horizontal: "right" }}
          >
            <Alert
              onClose={handleClose}
              severity="success"
              sx={{
                width: "100%",
                fontSize: "1.2em",
                padding: "20px",
                maxWidth: "600px",
              }}
            >
              Comentario realizado con exito!
            </Alert>
          </Snackbar>
          <Snackbar
            open={comentErr}
            autoHideDuration={6000}
            onClose={handleClose}
            anchorOrigin={{ vertical: "top", horizontal: "right" }}
          >
            <Alert
              onClose={handleClose}
              severity="error"
              sx={{
                width: "100%",
                fontSize: "1.2em",
                padding: "20px",
                maxWidth: "600px",
              }}
            >
              ¡No se aceptan comentarios vacios!
            </Alert>
          </Snackbar>
        </Paper>
        <Grid container spacing={2} alignItems="center">
          {/* Seccion de archivos*/}
          <Grid item xs={12} sx={{ marginTop: "5px" }}>
            <Files />
          </Grid>

          {/* Formulario para agregar comentario */}
          {!subscripto ? (
            ""
          ) : (
            <Grid item xs={12}>
              <Paper
                elevation={3}
                sx={{
                  padding: "10px",
                  backgroundColor: "#f0f0f0",
                  borderRadius: 1,
                  marginTop: "20px",
                  border: "3px solid #785589",
                }}
              >
                <Typography variant="h5" gutterBottom>
                  {" "}
                  Dejanos tu opinion del curso:
                </Typography>
                <InputBase
                  fullWidth
                  placeholder="Agregar comentario"
                  sx={{ paddingLeft: 2, maxWidth: "900px" }}
                  inputProps={{
                    "aria-label": "agregar comentario",
                    value: comentarioText,
                    onChange: (e) => setComentarioText(e.target.value),
                  }}
                />
                <Rating
                  name="simple-controlled"
                  sx={{ paddingLeft: 2 }}
                  value={value}
                  onChange={(event, newValue) => {
                    setValue(newValue);
                  }}
                />
                <br></br>
                <Button
                  variant="contained"
                  className="button-comentar"
                  value
                  sx={{ marginLeft: 2, marginTop: 2 }}
                  onClick={handleComment}
                >
                  Enviar
                </Button>
              </Paper>
            </Grid>
          )}

          {/* Sección de comentarios */}
          <Grid item xs={12}>
            <Comments />
          </Grid>
        </Grid>
      </Container>
    </>
  );
}

export default Course;
