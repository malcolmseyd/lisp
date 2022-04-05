(define list (lambda (. rest) rest))

(defmacro not (a)
  (if a nil #t))

(define map
  (lambda (f ls)
    (if ls
        (cons (f (car ls)) (map f (cdr ls)))
        ls)))

(define filter
  (lambda (f ls)
    (if ls
        (if (f (car ls))
            (cons (car ls) (filter f (cdr ls)))
            (filter f (cdr ls)))
        ls)))
