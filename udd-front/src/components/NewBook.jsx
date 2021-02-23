import React, { useState } from 'react'
import axios from 'axios'
import { Button, Container, Form, Grid } from 'semantic-ui-react'
import { useHistory } from 'react-router-dom'
import Moment from 'react-moment';

const NewBook = () => {
    const [author, setAuthor] = useState('')
    const [isbn, setIsbn] = useState('')
    const [title, setTitle] = useState('')
    const [genre,setGenre] = useState('')
    const [publishYear, setPublishYear] = useState('')
    const [fileName,setFilename] = useState('')
    const [file,setFile] = useState({})

    const history = useHistory()

    const handleSubmit = (event) => {
        event.preventDefault()
        console.log('fileName' + fileName)
        console.log('file' + file)
        
        const formData = new FormData()
        let newFileName = title
        var ext = fileName.split(".").pop()
        formData.append('file',file, newFileName+`.${ext}`)
        console.log(new Date(publishYear).toUTCString())// Date.parse(publishYear)
        axios.post('http://localhost:8080/upload', formData, { headers:{'content-type':'multipart/form-data'}})
            .then(resp=> {
                console.log(resp.data)
                // Here send another post request to send book data
                axios.post('http://localhost:8080/book', {
                    author,
                    isbn,
                    title,
                    genre,
                    publish_year: new Date(publishYear).toISOString(),
                },{
                    headers: {'content-type':"application/json"},
                })
                .then(resp => {
                    console.log(resp.data)
                    history.push('/books')
                })
                .catch(err => {console.log(err)})
            })
            .catch(e => {
                console.log(e)
            })
    }

    return (
        <Container>
            <Grid>
                <Grid.Row centered>
                    <Grid.Column width={6}>
                    <Form>
                            <Form.Field>
                                <label>Author</label>
                                <input placeholder='Author' autoComplete="off" value={author} onChange={(e)=>setAuthor(e.target.value)} />
                            </Form.Field>
                            <Form.Field>
                                <label>ISBN</label>
                                <input placeholder='Isbn' type="text" autoComplete="off" value={isbn} onChange={(e)=>setIsbn(e.target.value)}/>
                            </Form.Field>
                            <Form.Field>
                                <label>Title</label>
                                <input placeholder='title' autoComplete="off" value={title} onChange={(e)=>setTitle(e.target.value)}/>
                            </Form.Field>
                            <Form.Field>
                                <label>Genre</label>
                                <select placeholder="Genre"  autoComplete="off" value={genre} onChange={(e)=>setGenre(e.target.value)}>
                                    <option value="Comedy">Comedy</option>
                                    <option value="Horror">Horror</option>
                                    <option value="Thriller">Thriller</option>
                                    <option value="Drama">Drama</option>
                                    <option value="Romantic">Romantic</option>
                                </select>
                            </Form.Field>
                            <Form.Field>
                                <label>Publish Book</label>
                                <input placeholder="Publish year" type="date"  autoComplete="off" value={publishYear} onChange={(e)=>setPublishYear(e.target.value)}/>
                            </Form.Field>
                            <Form.Field>
                                <label>Book</label>
                                <input type="file" placeholder="Description"  autoComplete="off" onChange={(e) => { 
                                    setFile(e.target.files[0])
                                    setFilename(e.target.value)
                                    }}  />
                            </Form.Field>
                            <Button color="green" onClick={handleSubmit}>Submit</Button>
                        
                    </Form>
                    </Grid.Column>
                </Grid.Row>
            </Grid>
            <br/>
            <div style={{
                alignContent: 'center',
                backgroundImage: "url(/bck.jpg)",
                backgroundSize: "contain",
                backgroundRepeat: 'no-repeat',
                width: 1200,
                height: 800,
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
                paddingtop: "66%",
            }}>

            </div>
        </Container>
    )
}

export default NewBook
