import React, { useEffect, useState } from "react";
import { Paper, Typography, Button } from "@mui/material";

const Files = () => {
  const courseID = localStorage.getItem("courseID");
  const [archivos, setArchivos] = useState([]);
  const [showFiles, setShowFiles] = useState(true);
  const courseName = localStorage.getItem("cursoTitulo");

  useEffect(() => {
    if (showFiles) {
      getFiles();
    }
  }, [showFiles]);

  const getFiles = async () => {
    const response = await fetch(
      `http://localhost:8080/files/file/${courseID}`,
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      }
    );

    if (response.ok) {
      const data = await response.json();
      setArchivos(
        data.map((archivo) => ({
          ...archivo,
        }))
      );
    }
  };

  const downloadFile = (base64Content) => {
    try {
      // Decodificar base64 a datos binarios
      const byteCharacters = atob(base64Content);
      const byteNumbers = new Array(byteCharacters.length);
      for (let i = 0; i < byteCharacters.length; i++) {
        byteNumbers[i] = byteCharacters.charCodeAt(i);
      }
      const byteArray = new Uint8Array(byteNumbers);

      // Crear Blob con el tipo correcto (en este caso, PDF)
      const blob = new Blob([byteArray], { type: "application/pdf" });

      // Crear URL del Blob
      const url = URL.createObjectURL(blob);

      // Crear un enlace invisible y hacer click en él para descargar el archivo
      const a = document.createElement("a");
      a.href = url;
      a.download = courseName; // Nombre del archivo para la descarga
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);

      // Liberar la URL del Blob cuando ya no se necesite
      URL.revokeObjectURL(url);
    } catch (error) {
      console.error("Error al descargar el archivo:", error);
    }
  };

  return (
    <Paper
      sx={{
        padding: "20px",
        marginTop: "20px",
        border: "3px solid #785589",
        marginBottom: "20px",
        backgroundColor: "#f0f0f0",
      }}
    >
      <Typography variant="h5" gutterBottom color="#785589">
        Archivos
      </Typography>
      {archivos.length === 0 ? (
        <Typography variant="body1">No hay archivos aún.</Typography>
      ) : (
        archivos.map((archivo, index) => (
          <div key={archivo.ID} style={{ marginTop: "20px" }}>
            <Typography
              variant="body1"
              sx={{ fontFamily: "Arial", fontSize: "1.2em" }}
            >
              {archivo.fileName}
            </Typography>
            <Button
              variant="outlined"
              onClick={() => downloadFile(archivo.file)}
            >
              Descargar
            </Button>
          </div>
        ))
      )}
    </Paper>
  );
};

export default Files;
