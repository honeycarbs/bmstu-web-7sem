import { useState, useEffect, FormEvent, useRef } from "react";
import { Button, Col, Container, Form, Row, Stack } from "react-bootstrap";
import { useNavigate, Link } from "react-router-dom";
import TagService from "../service/TagService";
import Tag from "../types/Tag";
import NoteForm from "./NoteForm";
import { OnCreateNote } from "./OnCreateNote";
import CreatableReactSelect from "react-select/creatable"
import NoteService from "../service/NoteService";
import { CustomSelectStyle } from "./SelectStyles";


const NewNote: React.FC = () => {
    const [AvailableTags, setAvailableTags] = useState<Tag[]>([])
    const headerRef = useRef<HTMLInputElement>(null)
    const markdownRef = useRef<HTMLTextAreaElement>(null)

    const navigate = useNavigate()

    const [selectedTags, setSelectedTags] = useState<Tag[]>([])

    function NoteSubmitHandler(e: FormEvent) {
        e.preventDefault()

        let data = {
            header: headerRef.current!.value,
            body: markdownRef.current!.value,
            tags: selectedTags
        }
        OnCreateNote(data)
        navigate("..")
        // window.location.reload()
    }

    useEffect(() => {
        FetchTags();
    }, []);

    const FetchTags = () => {
        TagService.getAll()
            .then((response: any) => {
                setAvailableTags(response.tags);
            })
            .catch((e: Error) => {
                console.log(e);
            });
    }

    return (
        <Container className="my-4">
            <Form onSubmit={NoteSubmitHandler}>
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
                                    <Form.Control
                                        placeholder="header"
                                        required
                                        autoComplete="off"
                                        ref={headerRef}
                                    />
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
                                />
                            </Form.Group>
                        </Row>
                    </Stack>
                </Row>
            </Form>
        </Container>
    )
}

export default NewNote