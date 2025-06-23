import React, { useState, useEffect } from "react";
import "./Home.css";
import {
  AppBar,
  Toolbar,
  Typography,
  Button,
  Avatar,
  IconButton,
  Menu,
  List,
  ListItemIcon,
  ListItemText,
  ListItemButton,
  Grid,
  Container,
} from "@mui/material";
import AccountCircleIcon from "@mui/icons-material/AccountCircle";
import LogoutIcon from "@mui/icons-material/Logout";
import { useNavigate } from "react-router-dom";
import "../components/Componentes.css";

function Microservicios() {
  const navigate = useNavigate();
  const [logoutOpen, setLogoutOpen] = useState(false);
  const [instances, setInstances] = useState([]); // Estado para las instancias
  //  const userId = localStorage.getItem("userID");

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

  const fetchInstances = async () => {
    try {
      const response = await fetch("http://localhost:8087/admin/services", {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (response.ok) {
        const data = await response.json(); // Parsear el cuerpo JSON
        setInstances(data.instances); // Asignar las instancias al estado
      } else {
        console.error("Error fetching instances:", response.statusText);
      }
    } catch (error) {
      console.error("Error fetching instances:", error);
    }
  };

  // Cargar las instancias al montar el componente
  useEffect(() => {
    fetchInstances();
  }, []);

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

      {/* Sección para mostrar las instancias */}
      <Container sx={{ marginTop: 4 }}>
        <Typography variant="h5" sx={{ marginBottom: 2 }}>
          Instancias de Microservicios
        </Typography>
        <Grid container spacing={2}>
          {instances.length > 0 ? (
            instances.map((instance, index) => (
              <Grid item xs={12} sm={6} md={4} key={index}>
                <Button variant="outlined" fullWidth sx={{ padding: 2 }}>
                  {instance.nombre}{" "}
                  {/* Cambia "name" a "nombre" si así viene del backend */}
                </Button>
              </Grid>
            ))
          ) : (
            <Typography>No hay instancias activas en este momento.</Typography>
          )}
        </Grid>
      </Container>
    </>
  );
}

export default Microservicios;
