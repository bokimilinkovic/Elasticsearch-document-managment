import React, { useEffect, useState } from 'react'
import axios from 'axios'
import { Button, Form, Header, Image, Table } from 'semantic-ui-react'

const Users = () => {
    const [users,setUsers] = useState([])
    const [startingPoint, setStartingPoint] = useState('Novi Sad')
    const [range,setRange] = useState(0)
    
    useEffect(() => {
        axios.get(`${process.env.REACT_APP_API_URL}/users`,{
            headers: {
                'content-type':'aplication/json'
            }
        }).then(resp => {
            setUsers(resp.data)
        }).catch(error => {
            console.log(error)
        })
        console.log(users)
    }, [])

    const handleChange = (event) => {
        setStartingPoint(event.target.value)
    }

    const handleSubmit = (event) => {
        event.preventDefault()
        console.log(range + '' + startingPoint)
        axios.post(`${process.env.REACT_APP_API_URL}/distance`,{
            city: startingPoint,
            range: parseInt(range)
        },{
            headers: {'content-type':'application/json'}
        }).then(resp => {
            setUsers(resp.data)
        }).catch(error=>{
            console.log(error)
        })

    }

    return (
        <div>
            <Form size={"small"} >
                <Form.Field>
                    <label>Select your starting point</label>
                    <select value={startingPoint} onChange={handleChange}>
                        <option value="Novi Sad" selected>Novi Sad</option>
                        <option value="Belgrade">Belgrade</option>
                    </select>
                    <label>Enter km range</label>
                    <input type="number" value={range} onChange={(e) => setRange(e.target.value)} />
                </Form.Field>
                <Button type='submit' onClick={handleSubmit}>Search</Button>
            </Form>
            <Table celled padded >
                <Table.Header>
                <Table.Row>
                    <Table.HeaderCell singleLine>Username</Table.HeaderCell>
                    <Table.HeaderCell>City</Table.HeaderCell>
                    <Table.HeaderCell>Country</Table.HeaderCell>
                </Table.Row>
                </Table.Header>

                <Table.Body>
                    {users.map((user)=>{
                        let random = Math.floor(Math.random() * 100);  
                        return (
                            <Table.Row key={user.username}>
                                <Table.Cell>
                                    <Header as='h3' textAlign='center'>
                                        {user.username}
                                    </Header>
                                </Table.Cell>
                                <Table.Cell singleLine>{user.city}</Table.Cell>
                                <Table.Cell>
                                    {user.country}
                                </Table.Cell>
                                <Table.Cell>
                                    <Image src={`https://randomuser.me/api/portraits/men/${random}.jpg`}></Image>
                                </Table.Cell>
                            </Table.Row>
                        )
                    })}
                
                </Table.Body>
            </Table>
        </div>
    )
}

export default Users
