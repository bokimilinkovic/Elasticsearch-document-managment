import axios from 'axios'
import React, { Fragment, useEffect, useState } from 'react'
import { Header, Table, Rating, Form, Button } from 'semantic-ui-react'

const Books = () => {
    const [books, setBooks] = useState([])
    const [term, setTerm] = useState('')

    const handleSubmit = (event) => {
        event.preventDefault()
        console.log(term)
        axios.get(`http://localhost:8080/search?term=${term}`, {
            headers: {'content-type':'application/json'}
        }).then(resp=>{
           setBooks(resp.data)
        }).catch(error => {
            console.log(error)
        })
    }

    useEffect(() => {
        axios.get(`${process.env.REACT_APP_API_URL}/books`,{
            headers: {
                'content-type':'aplication/json'
            }
        }).then(resp => {
            setBooks(resp.data)
        }).catch(error => {
            console.log(error)
        })
        console.log(books)
    }, [])

    return (
        <div>
            <Form size={"small"} >
                <Form.Field>
                    <label>Search book by title, author, isbn, content or anything!</label>
                    <input value={term}onChange={(e) => setTerm(e.target.value)} />
                </Form.Field>
                <Button type='submit' onClick={handleSubmit}>Search</Button>
            </Form>
            <Table celled padded inverted>
                <Table.Header>
                <Table.Row>
                    <Table.HeaderCell singleLine>Book title</Table.HeaderCell>
                    <Table.HeaderCell>Author</Table.HeaderCell>
                    <Table.HeaderCell>Genre</Table.HeaderCell>
                    <Table.HeaderCell>Publish Year</Table.HeaderCell>
                    <Table.HeaderCell width={10}>Synposis</Table.HeaderCell>
                </Table.Row>
                </Table.Header>

                <Table.Body>
                    {books.map((book)=>{
                        return (
                            <Table.Row key={book.title}>
                                <Table.Cell>
                                    <Header as='h2' textAlign='center'>
                                        {book.title}
                                    </Header>
                                </Table.Cell>
                                <Table.Cell singleLine>{book.author}</Table.Cell>
                                <Table.Cell>
                                    {book.genre ? book.genre : 'drama'}
                                </Table.Cell>
                                <Table.Cell>
                                    {book.publish_year.toString().substring(0,4)}
                                </Table.Cell>
                                <Table.Cell  >
                                    {book.content.substring(0,100)}
                                </Table.Cell>
                                <Table.Cell>
                                    <a href={'http://localhost:8080/static/' +book.title + '.pdf'}>Preview and download</a>
                                </Table.Cell>
                            </Table.Row>
                        )
                    })}
                
                </Table.Body>
            </Table>
        </div>
    )
}

export default Books
