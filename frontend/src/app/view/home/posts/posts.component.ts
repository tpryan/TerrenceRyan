import { Component, OnInit } from '@angular/core';
import { PostService, Post} from '../../../service/post/post.service';

@Component({
  selector: 'app-posts',
  templateUrl: './posts.component.html',
  styleUrls: ['./posts.component.scss']
})
export class PostsComponent implements OnInit {

  posts: Array<Post> = [];

  constructor(private postService: PostService) {} 

  getPosts(): void {
    this.postService.getPosts().subscribe(posts => this.posts =  posts.slice(-4, -1));
  } 

  ngOnInit(): void {
    this.getPosts();
  }

}
