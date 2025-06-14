import { Box, Button, Container, Paper, Typography, IconButton, useTheme, Stack, Tooltip, Snackbar, Alert } from "@mui/material";
import { useState, useEffect } from "react";
import Editor from "@monaco-editor/react";
import Split from "react-split";
import { Brightness4, Brightness7, Lightbulb } from "@mui/icons-material";

export default function MainScreen({
  level,
  setLevel,
  code,
  setCode,
  output,
  user,
  handleRunCode,
  handleHint,
  hint,
  darkMode,
  toggleTheme,
  setShowProfile,
  levelData,
  isLoading
}) {
  const [openSuccess, setOpenSuccess] = useState(false);
  const theme = useTheme();
  const levelsNum = 16

  useEffect(() => {
    // Показываем уведомление, если уровень выполнен и это не нулевой уровень
    if (level > 0 && level <= user.completedLevels) {
      setOpenSuccess(true);
    }
  }, [level, user.completedLevels]);

  const handleCloseSuccess = () => {
    setOpenSuccess(false);
  };

  return (
    <Container maxWidth="xl" disableGutters sx={{ height: '100vh', display: 'flex', flexDirection: 'column' }}>
      <Snackbar
        open={openSuccess}
        autoHideDuration={6000}
        onClose={handleCloseSuccess}
        anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
      >
        <Alert 
          onClose={handleCloseSuccess} 
          severity="success"
          sx={{ width: '100%' }}
        >
          Уровень {level} успешно пройден! Можете перейти к следующему.
        </Alert>
      </Snackbar>
      <Box sx={{ p: 2 }}>
                <Paper elevation={3} sx={{ p: 2, borderRadius: 2 }}>
                    <Box sx={{ display: 'flex', flexWrap: 'nowrap', gap: 1, justifyContent:"flex-start", overflowX: "auto", 
                    '&::-webkit-scrollbar': {
                      height: '6px',
                    },
                    '&::-webkit-scrollbar-track': {
                      background: theme.palette.grey[200],
                      borderRadius: '3px',
                    },
                    '&::-webkit-scrollbar-thumb': {
                      background: theme.palette.primary.main,
                      borderRadius: '3px',
                    },
                    pb: 1,
                    }}>
                        {Array.from({ length: levelsNum }, (_, i) => {
                          const isCompleted = i <= user.completedLevels;
                          const isNext = i === user.completedLevels + 1;
                          const isCurrent = i === level;

                          return (
                            <Button 
                              key={i}
                              size="small"
                              variant={isCurrent ? "contained" : "outlined"}
                              color={isCompleted ? "success" : isNext ? "primary" : "inherit"}
                              onClick={() => setLevel(i)}
                              disabled={!isCompleted && !isNext}
                              sx={{ flex: '1 0 auto', m: 0.5,whiteSpace: 'nowrap', minWidth: '40px' }}
                            >
                              {i}
                            </Button>
                          );
                        })}
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
              <Typography variant="h6">Задача уровня {levelData.id}</Typography>
              
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
              <Typography sx={{ mb: 2, whiteSpace: 'pre-line'  }}>
                {isLoading ? "Загрузка уровня..." :levelData.description}
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
                  defaultLanguage={levelData.language}
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
              height: '100%',
              overflow: 'hidden'
            }}>
              <Typography variant="subtitle1" sx={{ mb: 1 }}>Результат выполнения:</Typography>
              <Box 
                sx={{
                  flex: 1,
                  border: `1px solid ${theme.palette.divider}`,
                  borderRadius: 1,
                  p: 1,
                  backgroundColor: theme.palette.background.paper,
                  overflow: 'auto',
                  color: output.includes('200') ? 
                    theme.palette.success.main : 
                    theme.palette.error.main,
                  fontFamily: 'monospace',
                  whiteSpace: 'pre-wrap',
                  wordBreak: 'break-word'
                }}
              >
                {output || <span style={{ color: theme.palette.text.secondary }}>Здесь будет отображаться результат выполнения кода</span>}
              </Box>
            </Paper>
          </Split>
        </Split>
      </Box>
    </Container>
  );
}