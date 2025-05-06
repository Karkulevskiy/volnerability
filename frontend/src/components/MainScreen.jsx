import { Box, Button, Container, Paper, TextField, Typography } from "@mui/material";
import Editor from "@monaco-editor/react";
import Split from "react-split";

export default function MainScreen({ level, setLevel, code, setCode, output, handleRunCode, handleHint, hint, darkMode }) {
  return (
    <Container maxWidth="xl" disableGutters sx={{ height: '100vh', display: 'flex', flexDirection: 'column' }}>
      <Box sx={{ p: 2 }}>
        <Paper elevation={3} sx={{ p: 2, borderRadius: 2 }}>
          <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1 }}>
            {[...Array(20)].map((_, index) => (
              <Button key={index} size="small" variant={level === index + 1 ? "contained" : "outlined"}
                color={level > index ? "success" : "inherit"} onClick={() => setLevel(index + 1)}>
                {index + 1}
              </Button>
            ))}
          </Box>
        </Paper>
      </Box>

      <Box sx={{ flex: 1, overflow: 'hidden', px: 2, pb: 2 }}>
        <Split className="split-horizontal" sizes={[30, 70]} minSize={250} gutterSize={10} direction="horizontal"
          style={{ display: 'flex', height: '100%' }}>
          <Paper elevation={3} sx={{ p: 2, borderRadius: 3, height: '100%', overflowY: 'auto' }}>
            <Typography variant="h6">Задача уровня {level}</Typography>
            <Typography>Найди уязвимость и получи доступ к секретным данным!</Typography>
            <Button onClick={handleHint} sx={{ mt: 1 }} variant="outlined">Получить подсказку</Button>
            {hint && <Typography sx={{ mt: 1 }} color="secondary">{hint}</Typography>}
          </Paper>

          <Split className="split-vertical" sizes={[70, 30]} minSize={100} gutterSize={8} direction="vertical"
            style={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
            <Paper sx={{ display: 'flex', flexDirection: 'column', flex: 1, p: 2, borderRadius: 3 }}>
              <Box sx={{ flexGrow: 1, minHeight: 0 }}>
                <Editor
                  defaultLanguage="python"
                  value={code}
                  onChange={(val) => setCode(val || "")}
                  theme={darkMode ? "vs-dark" : "light"}
                  options={{ fontSize: 18, minimap: { enabled: false }, lineNumbers: "on" }}
                  height="100%"
                />
              </Box>

              <Box sx={{ mt: 2, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <Button variant="contained" color="success" onClick={handleRunCode}>
                  Отправить код
                </Button>
              </Box>
            </Paper>
            <Paper elevation={3} sx={{ p: 2, borderRadius: 3 }}>
              <Typography variant="subtitle1">Результат выполнения:</Typography>
              <TextField fullWidth multiline rows={2} value={output} InputProps={{ readOnly: true }} />
            </Paper>
          </Split>
        </Split>
      </Box>
    </Container>
  );
}

