import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { TextField, Button, Typography, Paper, Grid, Container, Box, Snackbar, Alert, List, ListItem, ListItemText, ListItemSecondaryAction, IconButton } from '@mui/material';
import DeleteIcon from '@mui/icons-material/Delete';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';
import StopIcon from '@mui/icons-material/Stop';

const App = () => {
    const [vmName, setVmName] = useState('');
    const [vms, setVms] = useState([]);
    const [alert, setAlert] = useState({ open: false, message: '', severity: 'success' });

    useEffect(() => {
        fetchVMs();
    }, []);

    const fetchVMs = async () => {
        try {
            const response = await axios.get('http://localhost:8080/vms');
            setVms(Object.entries(response.data));
        } catch (error) {
            console.error(error);
            setAlert({ open: true, message: 'Failed to fetch VMs', severity: 'error' });
        }
    };

    const createVM = async () => {
        if (!vmName) return;
        try {
            await axios.post(`http://localhost:8080/create?name=${vmName}`);
            setVmName('');
            fetchVMs();
            setAlert({ open: true, message: `VM '${vmName}' created!`, severity: 'success' });
        } catch (error) {
            console.error(error);
            setAlert({ open: true, message: 'Failed to create VM', severity: 'error' });
        }
    };

    const deleteVM = async (name) => {
        try {
            await axios.delete(`http://localhost:8080/delete?name=${name}`);
            fetchVMs();
            setAlert({ open: true, message: `VM '${name}' deleted!`, severity: 'success' });
        } catch (error) {
            console.error(error);
            setAlert({ open: true, message: 'Failed to delete VM', severity: 'error' });
        }
    };

    const startVM = async (name) => {
        try {
            await axios.get(`http://localhost:8080/start?name=${name}`);
            fetchVMs();
            setAlert({ open: true, message: `VM '${name}' started!`, severity: 'success' });
        } catch (error) {
            console.error(error);
            setAlert({ open: true, message: 'Failed to start VM', severity: 'error' });
        }
    };

    const stopVM = async (name) => {
        try {
            await axios.get(`http://localhost:8080/stop?name=${name}`);
            fetchVMs();
            setAlert({ open: true, message: `VM '${name}' stopped!`, severity: 'success' });
        } catch (error) {
            console.error(error);
            setAlert({ open: true, message: 'Failed to stop VM', severity: 'error' });
        }
    };

    return (
        <Container maxWidth="sm">
            <Box my={4}>
                <Typography variant="h4" component="h1" gutterBottom>
                    VMware VM Manager
                </Typography>
                <TextField
                    label="VM Name"
                    variant="outlined"
                    value={vmName}
                    onChange={(e) => setVmName(e.target.value)}
                    fullWidth
                />
                <Button variant="contained" color="primary" onClick={createVM}>
                    Create VM
                </Button>
                <List>
                    {vms.map(([name]) => (
                        <ListItem key={name}>
                            <ListItemText primary={name} />
                            <ListItemSecondaryAction>
                                <IconButton edge="end" aria-label="start" onClick={() => startVM(name)}>
                                    <PlayArrowIcon />
                                </IconButton>
                                <IconButton edge="end" aria-label="stop" onClick={() => stopVM(name)}>
                                    <StopIcon />
                                </IconButton>
                                <IconButton edge="end" aria-label="delete" onClick={() => deleteVM(name)}>
                                    <DeleteIcon />
                                </IconButton>
                            </ListItemSecondaryAction>
                        </ListItem>
                    ))}
                </List>
            </Box>
            <Snackbar open={alert.open} autoHideDuration={6000} onClose={() => setAlert({ ...alert, open: false })}>
                <Alert onClose={() => setAlert({ ...alert, open: false })} severity={alert.severity}>
                    {alert.message}
                </Alert>
            </Snackbar>
        </Container>
    );
};

export default App;
