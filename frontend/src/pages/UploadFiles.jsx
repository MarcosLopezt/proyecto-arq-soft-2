import React, { useState } from "react";
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
  Box,
  Container,
  Snackbar,
  Alert,
} from "@mui/material";
import AccountCircleIcon from "@mui/icons-material/AccountCircle";
import LogoutIcon from "@mui/icons-material/Logout";
import { useNavigate } from "react-router-dom";
import "../components/Componentes.css";

function UploadFile() {
  const navigate = useNavigate();
  const [logoutOpen, setLogoutOpen] = useState(false);
  const courseID = parseInt(localStorage.getItem("courseID"), 10);
  const [selectedFile, setSelectedFile] = useState(null); // Estado para almacenar el archivo seleccionado
  const [base64String, setBase64String] = useState("");
  const [open, setOpen] = useState(false);
  //   const [wrongOpen, setWrongOpen] = useState(false);

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
  }; // Estado para almacenar el contenido en base64 del archivo

  //   const handleClose = (reason) => {
  //     if (reason === "clickaway") {
  //       return;
  //     }
  //     setOpen(false);
  //   };

  // Función para manejar el cambio en el input de archivo
  const handleFileInputChange = (event) => {
    const file = event.target.files[0];

    const reader = new FileReader();
    if (file) {
      reader.onloadend = () => {
        // Cuando la lectura del archivo esté completa, obtener el contenido en base64
        const base64data = reader.result;
        const base64Content = base64data.split(",")[1];

        setBase64String(base64Content);
        setSelectedFile(file); // Guardar el archivo seleccionado
      };

      reader.readAsDataURL(file); // Leer el archivo como base64
    }
  };

  const handleFileUpload = async () => {
    //console.log(base64String);
    const response = await fetch(`http://localhost:8080/files/upload`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ curso_id: courseID, file: base64String }),
    });

    console.log(response);

    if (response.ok) {
      //console.log("Se subio el archivo");
      setOpen(true);
    } else {
      //setWrongOpen(true);
      console.error("Error al subir el archivo:", response.statusText);
    }
  };

  const handleClose = (reason) => {
    if (reason === "clickaway") {
      return;
    }

    setOpen(false);
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
          <Button className="button-misCursos" variant="contained">
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

      <Container maxWidth="sm">
        <Box mt={4} p={3} boxShadow={3} bgcolor="background.paper">
          <Typography variant="h5" gutterBottom>
            Subir Archivo
          </Typography>
          <input
            type="file"
            onChange={handleFileInputChange}
            style={{ display: "none" }}
            id="file-upload"
          />
          <label htmlFor="file-upload">
            <Button variant="contained" component="span">
              Seleccionar Archivo
            </Button>
          </label>
          <Box mt={2}>
            {selectedFile && (
              <div>
                <Typography variant="subtitle1">
                  Nombre del Archivo: {selectedFile.name}
                </Typography>
                <Typography variant="subtitle1">
                  Tipo: {selectedFile.type}
                </Typography>
                <Typography variant="subtitle1">
                  Tamaño: {selectedFile.size} bytes
                </Typography>
              </div>
            )}
          </Box>
          <Box mt={2}>
            <Button
              variant="contained"
              color="primary"
              onClick={handleFileUpload}
              disabled={!selectedFile}
            >
              Subir
            </Button>
          </Box>
        </Box>

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
            Se subio el archivo con éxito
          </Alert>
        </Snackbar>
      </Container>
    </>
  );
}

export default UploadFile;
