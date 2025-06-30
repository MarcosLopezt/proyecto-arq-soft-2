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
  Card,
  CardContent,
  CardActions,
  Chip,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Alert,
  Snackbar,
} from "@mui/material";
import AccountCircleIcon from "@mui/icons-material/AccountCircle";
import LogoutIcon from "@mui/icons-material/Logout";
import AddIcon from "@mui/icons-material/Add";
import DeleteIcon from "@mui/icons-material/Delete";
import { useNavigate } from "react-router-dom";
import "../components/Componentes.css";

function Microservicios() {
  const navigate = useNavigate();
  const [logoutOpen, setLogoutOpen] = useState(false);
  const [instances, setInstances] = useState([]);
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [newInstance, setNewInstance] = useState({
    nombre: "",
    url: "",
  });
  const [snackbar, setSnackbar] = useState({
    open: false,
    message: "",
    severity: "success",
  });

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
        const data = await response.json();
        setInstances(data.instances);
      } else {
        console.error("Error fetching instances:", response.statusText);
        showSnackbar("Error al cargar las instancias", "error");
      }
    } catch (error) {
      console.error("Error fetching instances:", error);
      showSnackbar("Error de conexión al cargar las instancias", "error");
    }
  };

  const createInstance = async () => {
    try {
      const response = await fetch("http://localhost:8087/admin/services/create", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(newInstance),
      });

      if (response.ok) {
        showSnackbar("Instancia creada correctamente", "success");
        setCreateDialogOpen(false);
        setNewInstance({ nombre: "", url: "" });
        fetchInstances(); // Recargar la lista
      } else {
        const errorData = await response.json();
        showSnackbar(errorData.error || "Error al crear la instancia", "error");
      }
    } catch (error) {
      console.error("Error creating instance:", error);
      showSnackbar("Error de conexión al crear la instancia", "error");
    }
  };

  const deleteInstance = async (instanceId) => {
    try {
      const response = await fetch("http://localhost:8087/admin/services/remove", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ id: instanceId }),
      });

      if (response.ok) {
        showSnackbar("Instancia eliminada correctamente", "success");
        fetchInstances(); // Recargar la lista
      } else {
        const errorData = await response.json();
        showSnackbar(errorData.error || "Error al eliminar la instancia", "error");
      }
    } catch (error) {
      console.error("Error deleting instance:", error);
      showSnackbar("Error de conexión al eliminar la instancia", "error");
    }
  };

  const showSnackbar = (message, severity) => {
    setSnackbar({
      open: true,
      message,
      severity,
    });
  };

  const handleCloseSnackbar = () => {
    setSnackbar({ ...snackbar, open: false });
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
        <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: 2 }}>
          <Typography variant="h5">
            Gestión de Microservicios
          </Typography>
          <Button
            variant="contained"
            startIcon={<AddIcon />}
            onClick={() => setCreateDialogOpen(true)}
            sx={{ backgroundColor: "#785589" }}
          >
            Crear Instancia
          </Button>
        </div>
        
        <Grid container spacing={2}>
          {instances.length > 0 ? (
            instances.map((instance, index) => (
              <Grid item xs={12} sm={6} md={4} key={instance.id || index}>
                <Card sx={{ height: "100%" }}>
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      {instance.nombre}
                    </Typography>
                    <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                      ID: {instance.id}
                    </Typography>
                    {instance.url && (
                      <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                        URL: {instance.url}
                      </Typography>
                    )}
                    <Chip
                      label={instance.estado}
                      color={instance.estado === "Activo" ? "success" : "error"}
                      size="small"
                    />
                  </CardContent>
                  <CardActions>
                    {instance.id && instance.id.startsWith("inst_") && (
                      <Button
                        size="small"
                        color="error"
                        startIcon={<DeleteIcon />}
                        onClick={() => deleteInstance(instance.id)}
                      >
                        Eliminar
                      </Button>
                    )}
                  </CardActions>
                </Card>
              </Grid>
            ))
          ) : (
            <Grid item xs={12}>
              <Typography>No hay instancias activas en este momento.</Typography>
            </Grid>
          )}
        </Grid>
      </Container>

      {/* Dialog para crear nueva instancia */}
      <Dialog open={createDialogOpen} onClose={() => setCreateDialogOpen(false)}>
        <DialogTitle>Crear Nueva Instancia</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Nombre de la instancia"
            type="text"
            fullWidth
            variant="outlined"
            value={newInstance.nombre}
            onChange={(e) => setNewInstance({ ...newInstance, nombre: e.target.value })}
            sx={{ mb: 2 }}
          />
          <TextField
            margin="dense"
            label="URL de la instancia"
            type="url"
            fullWidth
            variant="outlined"
            value={newInstance.url}
            onChange={(e) => setNewInstance({ ...newInstance, url: e.target.value })}
            placeholder="http://ejemplo:puerto/ruta"
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setCreateDialogOpen(false)}>Cancelar</Button>
          <Button 
            onClick={createInstance}
            disabled={!newInstance.nombre || !newInstance.url}
            variant="contained"
          >
            Crear
          </Button>
        </DialogActions>
      </Dialog>

      {/* Snackbar para mostrar mensajes */}
      <Snackbar
        open={snackbar.open}
        autoHideDuration={6000}
        onClose={handleCloseSnackbar}
      >
        <Alert onClose={handleCloseSnackbar} severity={snackbar.severity}>
          {snackbar.message}
        </Alert>
      </Snackbar>
    </>
  );
}

export default Microservicios;
