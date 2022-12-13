import Tag from "../types/Tag"
import Note from "../types/Note"
import { useState, FormEvent, useEffect, useRef } from "react"
import { Row, Col, Stack, Button, Form, Container } from "react-bootstrap"
import { useNavigate, Link, useParams } from "react-router-dom"
import CreatableReactSelect from "react-select/creatable"
import { OnUpdateNote, UpdateNoteProps } from "./OnUpdateNote"
import TagService from "../service/TagService"
import NoteService from "../service/NoteService"
import { CustomSelectStyle } from "./SelectStyles"


const EditNote: React.FC = () => {
    const [AvailableTags, setAvailableTags] = useState<Tag[]>([])
    const [note, setNote] = useState<Note>({ id: null, header: "", body: "", tags: [] })

    const [selectedTags, setSelectedTags] = useState<Tag[]>(note.tags)
    const prevTags = note.tags

    const headerRef = useRef<HTMLInputElement>(null)
    const markdownRef = useRef<HTMLTextAreaElement>(null)

    const navigate = useNavigate()

    const params = useParams()


    function OnSubmitHandler(e: FormEvent) {
        e.preventDefault()

        let newNoteData: Note = {
            id: note.id,
            header: headerRef.current!.value,
            body: markdownRef.current!.value,
            tags: selectedTags,
        }


        OnUpdateNote(newNoteData, prevTags)

        navigate("..")
        window.location.reload()
    }

    useEffect(() => {
        retrieveTags();
        retrieveNote();
    }, []);

    const retrieveTags = () => {
        TagService.getAll()
            .then((response: any) => {
                setAvailableTags(response.tags);
            })
            .catch((e: Error) => {
                console.log(e);
            });
    };

    const retrieveNote = () => {
        NoteService.get(`/api/v1/notes/${params.id}`)
            .then((response: any) => {
                setNote({
                    id: response.id,
                    header: response.header,
                    body: response.body,
                    tags: response.tags
                });
                setSelectedTags(response.tags)
            })
            .catch((e: Error) => {
                console.log(e);
            });
    };


    return (
        <Container className="my-4">
            <Form onSubmit={OnSubmitHandler}>
                <Row className="align-items-center mb-3">
                    <Col>
                        <h2>my neat.ly</h2>
                    </Col>
                    <Col xs="auto">
                        <Stack direction="horizontal" gap={2} className="justify-content-end">
                            <Button type="submit">
                                Save
                            </Button>
                            <Link to="..">
                                <Button>
                                    Cancel
                                </Button>
                            </Link>
                        </Stack>
                    </Col>
                </Row>
                <Row>
                    <Stack gap={4}>
                        <Row>
                            <Col>
                                <Form.Group className="mb-3" controlId="header">
                                    <Form.Label>Header</Form.Label>
                                    <Form.Control placeholder="header" required
                                        autoComplete="off"
                                        ref={headerRef}
                                        defaultValue={note.header} />
                                </Form.Group>
                            </Col>

                            <Col>
                                <Form.Group className="mb-3" controlId="tags">
                                    <Form.Label>Tags</Form.Label>
                                    <CreatableReactSelect
                                        placeholder="tags"
                                        styles={CustomSelectStyle}
                                        value={selectedTags.map(tag => {
                                            return { label: tag.label, value: Number(tag.id) }
                                        })}
                                        onChange={tags => {
                                            setSelectedTags(
                                                tags.map(tag => {
                                                    return { label: tag.label, id: Number(tag.value) }
                                                })
                                            )
                                        }}
                                        onCreateOption={label => {
                                            const newTag = { id: 0, label: label }
                                            setSelectedTags(prev => [...prev, newTag])
                                        }}
                                        options={AvailableTags.map(tag => {
                                            return { label: tag.label, value: tag.id }
                                        })}
                                        isMulti />
                                </Form.Group>
                            </Col>
                            <Form.Group className="mb-3" controlId="markdown">
                                <Form.Label>Body</Form.Label>
                                <Form.Control placeholder="type here..."
                                    ref={markdownRef}
                                    as="textarea"
                                    rows={20}
                                    defaultValue={note.body} />
                            </Form.Group>
                        </Row>
                    </Stack>
                </Row>
            </Form>
        </Container>
    )
}

export default EditNote

