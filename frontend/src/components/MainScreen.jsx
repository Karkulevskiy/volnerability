import { Box, Button, Container, Paper, TextField, Typography, IconButton, useTheme, Stack, Tooltip } from "@mui/material";
import Editor from "@monaco-editor/react";
import Split from "react-split";
import { Brightness4, Brightness7, Lightbulb } from "@mui/icons-material";

export default function MainScreen({
  level,
  setLevel,
  code,
  setCode,
  output,
  handleRunCode,
  handleHint,
  hint,
  darkMode,
  toggleTheme,
  setShowProfile
}) {
  const theme = useTheme();

  return (
    <Container maxWidth="xl" disableGutters sx={{ height: '100vh', display: 'flex', flexDirection: 'column' }}>
      <Box sx={{ p: 2 }}>
        <Paper elevation={3} sx={{ p: 2, borderRadius: 2}}>
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
          <Paper elevation={3} sx={{ 
            p: 2, 
            borderRadius: 3, 
            height: '100%', 
            overflowY: 'auto',
            display: 'flex',
            flexDirection: 'column',
            position: 'relative'
          }}>
            {/* Основной контент задачи */}
            <Box sx={{ 
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
              mb: 2
            }}>
              <Typography variant="h6">Задача уровня {level}</Typography>
              
              <Tooltip title="Получить подсказку">
                <IconButton
                  onClick={handleHint}
                  color="secondary"
                  size="small"
                  sx={{
                    ml: 1,
                    '&:hover': {
                      backgroundColor: theme.palette.secondary.light,
                      color: theme.palette.secondary.contrastText
                    }
                  }}
                >
                  <Lightbulb fontSize="small" />
                </IconButton>
              </Tooltip>
            </Box>
            <Box sx={{ flexGrow: 1 }}>
              <Typography sx={{ mb: 2 }}>
                Найди уязвимость и получи доступ к секретным данным!
              </Typography>
              {hint && (
                <Typography sx={{ mt: 1 }} color="secondary">
                  {hint}
                </Typography>
              )}
            </Box>
            
              <Stack direction="row" spacing={1} sx={{
              mt: 'auto', // Прижимаем к низу
              pt: 2,      // Отступ сверху
              borderTop: `1px solid ${theme.palette.divider}`,
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
              flexShrink: 0 // Запрещаем уменьшение
            }}>
                <IconButton 
                  onClick={toggleTheme} 
                  color="inherit" 
                  size="medium"
                  sx={{ 
                    width: 44, 
                    height: 44,
                    border: `1px solid ${theme.palette.divider}`,
                    '&:hover': {
                      bgcolor: darkMode 
                        ? 'rgba(255, 255, 255, 0.08)' 
                        : 'rgba(0, 0, 0, 0.04)'
                    }
                  }}
                >
                  {darkMode ? <Brightness7 fontSize="medium" /> : <Brightness4 fontSize="medium" />}
                </IconButton>
                
                <Button 
                  variant="outlined" 
                  size="medium"
                  onClick={() => setShowProfile(true)}
                  sx={{ 
                    height: 44,
                    minWidth: 110,
                    textTransform: 'none',
                    px: 2
                  }}
                >
                  Профиль
                </Button>
              </Stack>
            
          </Paper>

          <Split className="split-vertical" 
            sizes={[70, 30]} 
            minSize={[150, 100]} // Минимальные размеры для обеих секций
            gutterSize={8} 
            direction="vertical"
            style={{ 
              display: 'flex', 
              flexDirection: 'column', 
              height: '100%',
              overflow: 'hidden' // Добавлено для корректного скролла
            }}>
            <Paper sx={{ display: 'flex', flexDirection: 'column', flex: 1, p: 2, borderRadius: 3, minHeight: 0  }}>
              <Box sx={{ flexGrow: 1, minHeight: 0, height: '100%', 
            overflowY: 'auto',
            display: 'flex',
            flexDirection: 'column' }}>
                <Editor
                  defaultLanguage="python"
                  value={code}
                  onChange={(val) => setCode(val || "")}
                  theme={darkMode ? "vs-dark" : "light"}
                  options={{ fontSize: 18, minimap: { enabled: false }, lineNumbers: "on", scrollBeyondLastLine: false }}
                  height="100%"
                  overflow="hidden"
                />
              </Box>

              <Box sx={{ mt: 2, flexShrink: 0, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <Button variant="contained" color="success" onClick={handleRunCode} sx={{
                      height: 40,
                      minWidth: 0
                    }}>
                  Отправить код
                </Button>
              </Box>
            </Paper>
            <Paper elevation={3} sx={{
              p: 2, 
              borderRadius: 3,
              display: 'flex',
              flexDirection: 'column',
              minHeight: 0, 
              overflow: 'hidden'}}>
              <Typography variant="subtitle1">Результат выполнения:</Typography>
              <TextField 
                fullWidth 
                multiline 
                rows={3} 
                value={output} 
                InputProps={{ 
                  readOnly: true,
                  sx: {
                    height: '100%',
                    '& textarea': {
                      overflow: 'auto !important',
                      resize: 'none'
                    }
                  } 
                }}
                sx={{
                  flex: 1,
                  minHeight: 0,
                  '& .MuiInputBase-root': {
                    height: '100%'
                  }
                }}
              />
            </Paper>
          </Split>
        </Split>
      </Box>
    </Container>
  );
}