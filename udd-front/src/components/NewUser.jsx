import React, { useState } from 'react'
import axios from 'axios'
import { Button, Container, Form, Grid } from 'semantic-ui-react'
import { useHistory } from 'react-router-dom'
import Moment from 'react-moment';

const NewUser = () => {
    const [username, setUsername] = useState('')
    const [city, setCity] = useState('')
    const [country, setCountry] = useState('')

    const history = useHistory()

    const handleSubmit = (event) => {
        event.preventDefault()
        axios.post('http://localhost:8080/users', {
            username,
            city,
            country,
        },{
            headers: {'content-type':"application/json"},
        })
        .then(resp => {
            console.log(resp.data)
            history.push('/users')
        })
        .catch(err => {console.log(err)})
    }
            

    return (
        <Container>
            <Grid>
                <Grid.Row centered>
                    <Grid.Column width={6}>
                    <Form>
                            <Form.Field>
                                <label>Username</label>
                                <input placeholder='username' autoComplete="off" value={username} onChange={(e)=>setUsername(e.target.value)} />
                            </Form.Field>
                            <Form.Field>
                                <label>City</label>
                                <input placeholder='city' type="text" autoComplete="off" value={city} onChange={(e)=>setCity(e.target.value)}/>
                            </Form.Field>
                            <Form.Field>
                                <label>Country</label>
                                <input placeholder='country' autoComplete="off" value={country} onChange={(e)=>setCountry(e.target.value)}/>
                            </Form.Field>
                            <Button color="green" onClick={handleSubmit}>Submit</Button>                    
                    </Form>
                    </Grid.Column>
                </Grid.Row>
            </Grid>
            <br/>
            <div style={{
                alignContent: 'center',
                backgroundImage: "url(/user.jpg)",
                backgroundSize: "contain",
                backgroundRepeat: 'no-repeat',
                width: 800,
                height: 600,
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
                paddingtop: "66%",
            }}>

            </div>
        </Container>
    )
}

export default NewUser